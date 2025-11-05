package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestParseCommitError(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected ErrorType
		message  string
	}{
		{
			name:     "hook failure",
			output:   "pre-commit hook failed",
			expected: ErrorTypeHookFailed,
			message:  "Pre-commit hook failed",
		},
		{
			name:     "no changes",
			output:   "nothing to commit",
			expected: ErrorTypeNoChanges,
			message:  "No changes to commit",
		},
		{
			name:     "merge conflict",
			output:   "fix conflicts then run git commit",
			expected: ErrorTypeMergeConflict,
			message:  "Merge conflicts need to be resolved",
		},
		{
			name:     "not in repo",
			output:   "not a git repository",
			expected: ErrorTypeNotInRepo,
			message:  "Not in a git repository",
		},
		{
			name:     "unknown error",
			output:   "some random error",
			expected: ErrorTypeUnknown,
			message:  "some random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseCommitError(nil, tt.output)
			if err.Type != tt.expected {
				t.Errorf("expected error type %v, got %v", tt.expected, err.Type)
			}
			if err.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, err.Message)
			}
		})
	}
}

func TestIsInGitRepo(t *testing.T) {
	// Test outside git repo
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)
	if IsInGitRepo() {
		t.Error("expected false outside git repo")
	}

	// Initialize git repo
	_, err := git.PlainInit(".", false)
	if err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	if !IsInGitRepo() {
		t.Error("expected true in git repo")
	}
}

func TestCommitSuccess(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo
	repo, err := git.PlainInit(".", false)
	if err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user using git commands
	cfg, err := repo.Config()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}
	cfg.User.Name = "Test User"
	cfg.User.Email = "test@example.com"
	err = repo.SetConfig(cfg)
	if err != nil {
		t.Fatalf("failed to set config: %v", err)
	}

	// Create and stage a file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Stage the file
	worktree, _ := repo.Worktree()
	_, err = worktree.Add("test.txt")
	if err != nil {
		t.Fatalf("failed to stage file: %v", err)
	}

	// Test successful commit
	err = Commit("feat: add test file")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCommitNoChanges(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo
	repo, err := git.PlainInit(".", false)
	if err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user
	cfg, err := repo.Config()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}
	cfg.User.Name = "Test User"
	cfg.User.Email = "test@example.com"
	err = repo.SetConfig(cfg)
	if err != nil {
		t.Fatalf("failed to set config: %v", err)
	}

	// Test commit with no changes
	err = Commit("feat: test commit")
	if err == nil {
		t.Error("expected error for no changes")
		return
	}

	commitErr, ok := err.(*CommitError)
	if !ok {
		t.Errorf("expected CommitError, got %T: %v", err, err)
		return
	}

	if commitErr.Type != ErrorTypeNoChanges {
		t.Errorf("expected error type %v, got %v. Output: %q", ErrorTypeNoChanges, commitErr.Type, commitErr.Output)
	}
}
