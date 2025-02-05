package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repository struct{}

func New() (*Repository, error) {
	if err := validateEnvironment(); err != nil {
		return nil, err
	}
	return &Repository{}, nil
}

// validateEnvironment checks if git is installed and we're in a git repository
func validateEnvironment() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found. please install git")
	}

	if err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		return fmt.Errorf("not in a git repository. please run this command in a git repository")
	}

	return nil
}

// GetCurrentBranch returns the name of the current branch
func (r *Repository) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	bytes, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return strings.TrimSpace(string(bytes)), nil
}
