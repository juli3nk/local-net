package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"dagger/local-net/internal/dagger"

	cplatforms "github.com/containerd/platforms"
	"github.com/juli3nk/go-utils/ci"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Build container images
func (m *LocalNet) Build(
	// +optional
	version string,
) (*LocalNet, error) {
	platformSpecifiers := []string{
		"linux/amd64",
	}
	platforms, err := cplatforms.ParseAll(platformSpecifiers)
	if err != nil {
		return nil, err
	}

	appVersion := ci.ResolveVersion(version, m.Git.Tag, m.Git.Commit, m.Git.Uncommitted)
	goAppVersionPkgPath := fmt.Sprintf("%s/pkg/version", appSourceUrl)
	tsNow := time.Now()

	goBuildPackages := []string{"."}
	goBuildLdflags := []string{
		fmt.Sprintf("-X %s.Version=%s", goAppVersionPkgPath, appVersion),
		fmt.Sprintf("-X %s.GitCommit=%s", goAppVersionPkgPath, m.Git.Commit),
		fmt.Sprintf("-X %s.BuildDate=%d", goAppVersionPkgPath, tsNow.Unix()),
	}

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(platforms))

	for _, platform := range platforms {
		wg.Add(1)
		go func(platform cplatforms.Platform) {
			defer wg.Done()
			adguardhomeBin := dag.Container().
				From(alpineBaseImage).
				WithExec([]string{"apk", "--update", "add",
					"curl",
					"tar",
				}).
				WithExec([]string{
					"curl",
					"-sfL",
					"-o",
					"/tmp/AdGuardHome.tar.gz",
					fmt.Sprintf("https://github.com/AdguardTeam/AdGuardHome/releases/download/%s/AdGuardHome_%s_%s.tar.gz", adguardhomeVersion, strings.ToLower(platform.OS), strings.ToLower(platform.Architecture))}).
				WithExec([]string{"tar", "-xz", "-f", "/tmp/AdGuardHome.tar.gz", "-C", "/tmp"}).
				File("/tmp/AdGuardHome/AdGuardHome")

			opts := dagger.GoBuildOpts{
				CgoEnabled: "1",
				Ldflags:    goBuildLdflags,
				Musl:       true,
				Arch:       platform.Architecture,
				Os:         platform.OS,
			}
			goBuilder := dag.Go(goVersion, m.Worktree).Build(appName, goBuildPackages, opts)

			binaryPath := fmt.Sprintf("/%s", appName)

			image := dag.Container(dagger.ContainerOpts{Platform: dagger.Platform(cplatforms.Format(platform))}).
				From(alpineBaseImage).
				WithExec([]string{"apk", "--update", "add", "networkmanager-cli"}).
				WithFile("/usr/local/bin/AdGuardHome", adguardhomeBin).
				WithFile(binaryPath, goBuilder).
				WithEntrypoint([]string{binaryPath}).
				WithLabel(specs.AnnotationCreated, tsNow.Format("2006-01-02T15:04:05 -0700")).
				WithLabel(specs.AnnotationSource, fmt.Sprintf("https://%s", appSourceUrl)).
				WithLabel(specs.AnnotationVersion, appVersion).
				WithLabel(specs.AnnotationRevision, m.Git.Commit).
				WithLabel(specs.AnnotationTitle, appName).
				WithLabel(specs.AnnotationDescription, "Local Net")

			m.Containers = append(m.Containers, image)
			errorsChan <- nil
		}(platform)
	}

	wg.Wait()
	close(errorsChan)

	var buildErrors []error
	for err := range errorsChan {
		if err != nil {
			buildErrors = append(buildErrors, err)
		}
	}

	if len(buildErrors) > 0 {
		return nil, fmt.Errorf("build failed: %w", errors.Join(buildErrors...))
	}

	return m, nil
}

func (m *LocalNet) Stdout(ctx context.Context) (string, error) {
	var outputs []string

	for _, ctr := range m.Containers {
		out, err := ctr.WithExec([]string{"cat", "/etc/os-release"}).Stdout(ctx)
		if err != nil {
			return "", err
		}

		outputs = append(outputs, out)
	}

	return strings.Join(outputs, "\n"), nil
}
