package main

import (
	"context"
	"fmt"

	"dagger/local-net/internal/dagger"
)

const (
	appName      = "local-net"
	appSourceUrl = "github.com/juli3nk/local-net"
)

type Git struct {
	Commit        string
	Tag           string
	Uncommitted   bool
	ModifiedFiles []string
}

type RegistryAuth struct {
	Address  string
	Username string
	Secret   *dagger.Secret
}

type LocalNet struct {
	Worktree     *dagger.Directory
	Git          *Git
	RegistryAuth *RegistryAuth
	Containers   []*dagger.Container
}

func fetchGitInfo(ctx context.Context, source *dagger.Directory) (*Git, error) {
	git := dag.Gitlocal(source)

	commit, err := git.GetLatestCommit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}
	tag, err := git.GetLatestTag(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest tag: %w", err)
	}
	uncommitted, err := git.Uncommitted(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check uncommitted changes: %w", err)
	}
	modifiedFiles, err := git.GetModifiedFiles(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	return &Git{
		Commit:        commit,
		Tag:           tag,
		Uncommitted:   uncommitted,
		ModifiedFiles: modifiedFiles,
	}, nil
}

func New(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	registryAddress string,
	// +optional
	registryUsername string,
	// +optional
	registrySecret *dagger.Secret,
) (*LocalNet, error) {
	git, err := fetchGitInfo(ctx, source)
	if err != nil {
		return nil, err
	}

	app := LocalNet{Worktree: source, Git: git}

	if len(registryAddress) > 0 {
		registryAuth := RegistryAuth{
			Address:  registryAddress,
			Username: registryUsername,
			Secret:   registrySecret,
		}

		app.RegistryAuth = &registryAuth
	}

	return &app, nil
}
