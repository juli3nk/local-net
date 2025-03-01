package main

import (
	"context"

	"dagger/local-net/internal/dagger"
)

// Lint commit messages
func (m *LocalNet) LintCommitMsg(
	ctx context.Context,
	args []string,
) (string, error) {
	return dag.Commitlint().
		Lint(m.Worktree, dagger.CommitlintLintOpts{Args: args}).
		Stdout(ctx)
}
