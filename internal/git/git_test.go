package git

import (
	"os"
	"os/exec"
	"testing"
)

func setupTestRepo(t *testing.T) (string, string, func()) {
	t.Helper()

	// Create temporary directory for git repo
	tmpDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create a separate temporary directory for non-git tests
	nonGitDir, err := os.MkdirTemp("", "non-git-test-*")
	if err != nil {
		t.Fatalf("failed to create non-git directory: %v", err)
	}

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Initialize git repo
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("failed to initialize git repo: %v", err)
	}

	// Configure git user for test environment
	if err := exec.Command("git", "config", "user.email", "test@example.com").Run(); err != nil {
		t.Fatalf("failed to configure git user email: %v", err)
	}
	if err := exec.Command("git", "config", "user.name", "Test User").Run(); err != nil {
		t.Fatalf("failed to configure git user name: %v", err)
	}

	// Create initial commit
	if err := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit").Run(); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	// Create test branches
	branches := []string{"test-branch-1", "test-branch-2"}
	for _, branch := range branches {
		if err := exec.Command("git", "checkout", "-b", branch).Run(); err != nil {
			t.Fatalf("failed to create branch %s: %v", branch, err)
		}
	}

	// Return to main branch
	if err := exec.Command("git", "checkout", "main").Run(); err != nil {
		// Try master for older git versions
		if err := exec.Command("git", "checkout", "master").Run(); err != nil {
			t.Fatalf("failed to return to main/master branch: %v", err)
		}
	}

	return tmpDir, nonGitDir, func() {
		os.Chdir(currentDir)
		os.RemoveAll(tmpDir)
		os.RemoveAll(nonGitDir)
	}
}

func TestNew(t *testing.T) {
	gitDir, nonGitDir, cleanup := setupTestRepo(t)
	defer cleanup()

	t.Run("inside git repository", func(t *testing.T) {
		// Ensure we're in the git repo directory
		if err := os.Chdir(gitDir); err != nil {
			t.Fatalf("failed to change to original directory: %v", err)
		}

		repo, err := New()
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
		} else if repo == nil {
			t.Errorf("New() returned nil, want a valid Repository instance")
		}
	})

	t.Run("outside git repository", func(t *testing.T) {
		// Change to the non-git directory
		if err := os.Chdir(nonGitDir); err != nil {
			t.Fatalf("failed to change to non-git directory: %v", err)
		}

		if _, err := New(); err == nil {
			t.Error("New() = nil, want error when not in git repository")
		}
	})
}

func TestGetCurrentBranch(t *testing.T) {
	gitDir, _, cleanup := setupTestRepo(t)
	defer cleanup()

	// create repository
	repo, err := New()
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// move to git directory
	if err := os.Chdir(gitDir); err != nil {
		t.Fatalf("failed to change to git directory: %v", err)
	}

	// git switch test-branch-1
	if err := exec.Command("git", "switch", "test-branch-1").Run(); err != nil {
		t.Fatalf("failed to switch to test-branch-1: %v", err)
	}

	t.Run("get current branch", func(t *testing.T) {
		branch, err := repo.GetCurrentBranch()

		if err != nil {
			t.Fatalf("GetCurrentBranch() error = %v, want nil", err)
		}

		if branch == "" {
			t.Fatalf("GetCurrentBranch() returned empty string, want branch name")
		}

		if branch != "test-branch-1" {
			t.Fatalf("GetCurrentBranch() returned %s, want test-branch-1", branch)
		}
	})
}
