package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tsukinoko-kun/request-review/internal/forge"
)

type (
	RepoInfo struct {
		RemoteURL string
		Branch    string
	}
)

func SmartPatch() (string, error) {
	cmd := exec.Command("git", "fetch", "origin", "main")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to fetch origin/main: %w", err)
	}

	cmd = exec.Command("git", "rev-parse", "origin/main")
	remoteHashB, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get origin/main HEAD: %w", err)
	}
	remoteHash := strings.TrimSpace(string(remoteHashB))

	cmd = exec.Command("git", "log", "-1", "--format=%H")
	cmd.Stdin = os.Stdin
	localHashB, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get local current commit: %w", err)
	}
	localHash := strings.TrimSpace(string(localHashB))

	if localHash == remoteHash {
		return "", errors.New("local and remote are up to date")
	}

	return Patch(remoteHash, localHash)
}

func Patch(from, to string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	r, err := git.PlainOpen(wd)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	fromHash, err := r.ResolveRevision(plumbing.Revision(from))
	if err != nil {
		return "", fmt.Errorf("failed to resolve revision %s: %w", from, err)
	}
	commit1, err := r.CommitObject(*fromHash)
	if err != nil {
		return "", fmt.Errorf("failed to get commit1: %w", err)
	}

	toHash, err := r.ResolveRevision(plumbing.Revision(to))
	if err != nil {
		return "", fmt.Errorf("failed to resolve revision %s: %w", to, err)
	}
	commit2, err := r.CommitObject(*toHash)
	if err != nil {
		return "", fmt.Errorf("failed to get commit2: %w", err)
	}

	patch, err := commit1.Patch(commit2)
	if err != nil {
		return "", fmt.Errorf("failed to generate patch: %w", err)
	}

	return patch.String(), nil
}

func RepoUrl() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetRepoInfo() (forge.ForgeInfo, error) {
	wd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, fmt.Errorf("failed to get working directory: %w", err)
	}
	r, err := git.PlainOpen(wd)
	if err != nil {
		return RepoInfo{}, fmt.Errorf("failed to open repository: %w", err)
	}

	origin, err := r.Remote("origin")
	if err != nil {
		return RepoInfo{}, fmt.Errorf("failed to get remote: %w", err)
	}

	remoteURL := origin.Config().URLs[0]

	branch, err := r.Head()
	if err != nil {
		return RepoInfo{}, fmt.Errorf("failed to get branch: %w", err)
	}

	branchName := branch.Name().Short()

	return RepoInfo{
		RemoteURL: remoteURL,
		Branch:    branchName,
	}, nil
}

func (ri RepoInfo) Name() string {
	pathParts := strings.Split(ri.RemoteURL, "/")
	name := pathParts[len(pathParts)-1]
	name = strings.TrimSuffix(name, ".git")
	return name
}

func (ri RepoInfo) Bookmark() string {
	return ri.Branch
}

func User() string {
	cmd := exec.Command("git", "config", "--global", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown User"
	}
	return strings.TrimSpace(string(output))
}
