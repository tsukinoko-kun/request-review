package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tsukinoko-kun/request-review/internal/config"
)

func Patch(cfg config.Config, from, to string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	r, err := git.PlainOpen(wd)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	commit1, err := r.CommitObject(plumbing.NewHash(from))
	if err != nil {
		return "", fmt.Errorf("failed to get commit1: %w", err)
	}

	commit2, err := r.CommitObject(plumbing.NewHash(to))
	if err != nil {
		return "", fmt.Errorf("failed to get commit2: %w", err)
	}

	patch, err := commit1.Patch(commit2)
	if err != nil {
		return "", fmt.Errorf("failed to generate patch: %w", err)
	}

	return patch.String(), nil
}
