// Package git provides git operations for the git-cc application with proper error handling.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type CommitError struct {
	Type    ErrorType
	Message string
	Output  string
	Err     error
}

type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeHookFailed
	ErrorTypeNoChanges
	ErrorTypeMergeConflict
	ErrorTypeNotInRepo
)

func (e *CommitError) Error() string {
	return fmt.Sprintf("git commit failed: %s", e.Message)
}

func (e *CommitError) IsHook() bool {
	return e.Type == ErrorTypeHookFailed
}

func (e *CommitError) GetDetails() string {
	if e.Output != "" {
		return e.Output
	}
	return e.Message
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)

	var outBuffer, errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		output := outBuffer.String() + errBuffer.String()
		return parseCommitError(err, output)
	}

	return nil
}

func parseCommitError(err error, output string) *CommitError {
	commitErr := &CommitError{
		Err:     err,
		Output:  output,
		Type:    ErrorTypeUnknown,
		Message: "Commit failed",
	}

	outputStr := strings.ToLower(output)

	switch {
	case strings.Contains(outputStr, "hook"):
		commitErr.Type = ErrorTypeHookFailed
		commitErr.Message = "Pre-commit hook failed"

	case strings.Contains(outputStr, "nothing to commit"):
		commitErr.Type = ErrorTypeNoChanges
		commitErr.Message = "No changes to commit"

	case strings.Contains(outputStr, "merge conflict") || strings.Contains(outputStr, "conflicts then run git commit"):
		commitErr.Type = ErrorTypeMergeConflict
		commitErr.Message = "Merge conflicts need to be resolved"

	case strings.Contains(outputStr, "not a git repository"):
		commitErr.Type = ErrorTypeNotInRepo
		commitErr.Message = "Not in a git repository"

	default:
		if output != "" {
			commitErr.Message = strings.TrimSpace(strings.Split(output, "\n")[0])
		}
	}

	return commitErr
}

func IsInGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// IsGitRepository is an alias for IsInGitRepo to match the main.go interface
func IsGitRepository() bool {
	return IsInGitRepo()
}

type CommitResult struct {
	Success bool
	Message string
	Details string
}

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	output := outBuffer.String()
	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	files := strings.Split(strings.TrimSpace(output), "\n")
	return files, nil
}

func CommitWithResult(message string) *CommitResult {
	err := Commit(message)
	if err != nil {
		if commitErr, ok := err.(*CommitError); ok {
			return &CommitResult{
				Success: false,
				Message: commitErr.Message,
				Details: commitErr.GetDetails(),
			}
		}
		return &CommitResult{
			Success: false,
			Message: "Commit failed",
			Details: err.Error(),
		}
	}

	return &CommitResult{
		Success: true,
		Message: "Changes committed successfully",
	}
}
