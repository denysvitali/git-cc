package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/denysvitali/git-cc/pkg/git"
	"github.com/denysvitali/git-cc/ui"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		printVersion()
		return
	}

	// Check if we're in a git repository
	if !git.IsGitRepository() {
		fmt.Printf("Error: not a git repository (or any of the parent directories): .git")
		os.Exit(1)
	}

	// Check if there are staged files
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		fmt.Printf("Error checking git status: %v", err)
		os.Exit(1)
	}

	if len(stagedFiles) == 0 {
		fmt.Printf("No staged files found. Stage files with 'git add' first.")
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("git-cc %s\n", version)
	fmt.Printf("  Commit: %s\n", commit)
	fmt.Printf("  Built: %s\n", date)
	if builtBy != "unknown" {
		fmt.Printf("  Built by: %s\n", builtBy)
	}
}
