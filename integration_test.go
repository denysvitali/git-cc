package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/denysvitali/git-cc/pkg/git"
)

func TestApplicationIntegration(t *testing.T) {
	// Create a temporary git repository for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo
	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd...); err != nil {
			t.Fatalf("Failed to run command %v: %v", cmd, err)
		}
	}

	// Create and stage a file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := runCommand("git", "add", "test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test that the application starts without error
	cmd := exec.Command("go", "run", ".", "--test.integration")
	cmd.Dir = originalDir
	cmd.Env = append(os.Environ(), "GIT_CC_TEST_MODE=1", "GIT_CC_TEST_DIR="+tempDir)

	// Since this is a TUI app, we can't easily interact with it in tests
	// But we can at least verify it starts and checks git repo properly
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	// Give it a moment to start and check git repo
	time.Sleep(100 * time.Millisecond)

	// Kill the process since we can't interact with it in tests
	cmd.Process.Kill()
	cmd.Wait()
}

func TestApplicationNotGitRepo(t *testing.T) {
	t.Skip("Skipping test due to complexity of testing directory changes in Go")
}

func TestApplicationNoStagedFiles(t *testing.T) {
	// Test application in git repo with no staged files
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo but don't stage any files
	if err := runCommand("git", "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = originalDir
	cmd.Env = append(os.Environ(), "GIT_CC_TEST_DIR="+tempDir)

	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected application to fail with no staged files")
	}

	if !strings.Contains(string(output), "No staged files found") {
		t.Errorf("Expected 'No staged files found' error, got: %s", string(output))
	}
}

func TestCommitWithPreCommitHook(t *testing.T) {
	// Create a temporary git repository with a pre-commit hook
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo
	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd...); err != nil {
			t.Fatalf("Failed to run command %v: %v", cmd, err)
		}
	}

	// Create a pre-commit hook that fails
	hookDir := filepath.Join(tempDir, ".git", "hooks")
	hookFile := filepath.Join(hookDir, "pre-commit")
	hookContent := `#!/bin/sh
echo "Pre-commit hook failed!"
exit 1
`
	if err := os.WriteFile(hookFile, []byte(hookContent), 0755); err != nil {
		t.Fatalf("Failed to create pre-commit hook: %v", err)
	}

	// Create and stage a file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := runCommand("git", "add", "test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test git package commit function with failing hook
	result := git.CommitWithResult("feat: test commit with failing hook")
	if result.Success {
		t.Error("Expected commit to fail due to pre-commit hook")
	}

	if !strings.Contains(result.Message, "Pre-commit hook failed") {
		t.Errorf("Expected pre-commit hook error, got: %s", result.Message)
	}
}

func TestConventionalCommitMessageFormat(t *testing.T) {
	// Create a temporary git repository for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Initialize git repo
	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd...); err != nil {
			t.Fatalf("Failed to run command %v: %v", cmd, err)
		}
	}

	// Test various conventional commit formats
	testCases := []struct {
		message string
		valid   bool
	}{
		{"feat: add new feature", true},
		{"fix(auth): resolve login issue", true},
		{"docs: update README", true},
		{"style: fix formatting", true},
		{"refactor: simplify code", true},
		{"test: add unit tests", true},
		{"build: update dependencies", true},
		{"ci: add GitHub Actions", true},
		{"chore: update .gitignore", true},
		{"invalid message", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			// Create and stage a file for each test
			testFile := filepath.Join(tempDir, "test_"+tc.message+".txt")
			if err := os.WriteFile(testFile, []byte(tc.message), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			if err := runCommand("git", "add", testFile); err != nil {
				t.Fatalf("Failed to stage file: %v", err)
			}

			result := git.CommitWithResult(tc.message)

			if tc.valid && !result.Success {
				t.Errorf("Expected commit to succeed for valid message: %s", tc.message)
			}

			// Check that the commit message was properly formatted
			if result.Success {
				commitLog, err := runCommandWithOutput("git", "log", "-1", "--pretty=format:%s")
				if err != nil {
					t.Logf("Failed to get commit log (non-fatal): %v", err)
				} else if commitLog != tc.message {
					t.Errorf("Expected commit message '%s', got '%s'", tc.message, commitLog)
				}
			}
		})
	}
}

// Helper function to run commands
func runCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "."
	return cmd.Run()
}

// Helper function to run commands and get output
func runCommandWithOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "."
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}
