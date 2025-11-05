// Package ui provides the terminal user interface for git-cc using Bubble Tea.
package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/denysvitali/git-cc/pkg/git"
)

type item struct {
	commitType  string
	description string
}

func (i item) Title() string       { return i.commitType }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.commitType + " " + i.description }

type Model struct {
	list      list.Model
	scope     textinput.Model
	message   textinput.Model
	step      int
	gitResult *git.CommitResult
	showError bool
}

const (
	StepTypeSelect = iota
	StepScope
	StepMessage
	StepError
)

const (
	paddingVertical   = 1
	paddingHorizontal = 2
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	appStyle = lipgloss.NewStyle().
			Padding(paddingVertical, paddingHorizontal)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	listItemStyle     = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type itemListDelegate struct{}

func (i itemListDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	if li == nil {
		return
	}

	itm := li.(item)
	style := listItemStyle
	if index == m.Cursor() {
		style = selectedItemStyle
		_, _ = io.WriteString(w, style.Render("❯ "))
	} else {
		_, _ = io.WriteString(w, "  ")
	}
	_, _ = io.WriteString(w, style.Render(fmt.Sprintf("%-10s %s", itm.commitType, itm.description)))
}

func (i itemListDelegate) Height() int {
	return 1
}

func (i itemListDelegate) Spacing() int {
	return 0
}

func (i itemListDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

var _ list.ItemDelegate = itemListDelegate{}

func InitialModel() Model {
	items := []list.Item{
		item{commitType: "feat", description: "A new feature"},
		item{commitType: "fix", description: "A bug fix"},
		item{commitType: "docs", description: "Documentation only changes"},
		item{commitType: "style", description: "Changes that do not affect the meaning of the code"},
		item{commitType: "refactor", description: "A code change that neither fixes a bug nor adds a feature"},
		item{commitType: "perf", description: "A code change that improves performance"},
		item{commitType: "test", description: "Adding missing tests or correcting existing tests"},
		item{commitType: "build", description: "Changes that affect the build system or external dependencies"},
		item{commitType: "ci", description: "Changes to CI configuration files and scripts"},
		item{commitType: "chore", description: "Other changes that don't modify src or test files"},
	}

	delegate := itemListDelegate{}
	commitList := list.New(items, delegate, 0, 0)
	commitList.Title = "Select the type of change"
	commitList.SetFilteringEnabled(true)
	commitList.SetShowHelp(true)

	scopeInput := textinput.New()
	scopeInput.Placeholder = "scope (optional)"
	scopeInput.CharLimit = 50
	scopeInput.Width = 30

	messageInput := textinput.New()
	messageInput.Placeholder = "commit message"
	messageInput.CharLimit = 100
	messageInput.Width = 50

	return Model{
		list:      commitList,
		scope:     scopeInput,
		message:   messageInput,
		step:      StepTypeSelect,
		showError: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch m.step {
			case StepTypeSelect:
				m.step = StepScope
				m.scope.Focus()
				return m, textinput.Blink

			case StepScope:
				m.step = StepMessage
				m.scope.Blur()
				m.message.Focus()
				return m, textinput.Blink

			case StepMessage:
				if m.message.Value() == "" {
					return m, nil
				}

				commitMsg := m.buildCommitMessage()
				m.gitResult = git.CommitWithResult(commitMsg)

				if !m.gitResult.Success {
					m.step = StepError
					m.showError = true
					return m, nil
				}

				return m, tea.Quit
			}

		case "r":
			if m.step == StepError && m.showError {
				// Retry - go back to message input
				m.step = StepMessage
				m.message.Focus()
				m.showError = false
				return m, textinput.Blink
			}
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	switch m.step {
	case StepTypeSelect:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

	case StepScope:
		m.scope, cmd = m.scope.Update(msg)
		cmds = append(cmds, cmd)

	case StepMessage:
		m.message, cmd = m.message.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string

	switch m.step {
	case StepTypeSelect:
		s = m.list.View()

	case StepScope:
		s = titleStyle.Render("Enter scope (optional, press Enter to skip):") + "\n"
		s += m.scope.View()

	case StepMessage:
		selectedItem := m.list.SelectedItem().(item)
		scopeStr := ""
		if m.scope.Value() != "" {
			scopeStr = fmt.Sprintf("(%s)", m.scope.Value())
		}
		s = titleStyle.Render("Enter commit message:") + "\n"
		s += promptStyle.Render(fmt.Sprintf("%s%s: ", selectedItem.commitType, scopeStr))
		s += m.message.View()

	case StepError:
		s = errorStyle.Render("Commit Failed!") + "\n\n"
		if m.gitResult != nil {
			s += m.gitResult.Message + "\n\n"
			if m.gitResult.Details != "" && m.gitResult.Details != m.gitResult.Message {
				s += m.gitResult.Details + "\n\n"
			}
		}
		s += promptStyle.Render("Press 'r' to retry or 'q' to quit")
	}

	return appStyle.Render(s)
}

func (m Model) buildCommitMessage() string {
	selectedItem := m.list.SelectedItem().(item)
	scopeStr := ""
	if m.scope.Value() != "" {
		scopeStr = fmt.Sprintf("(%s)", m.scope.Value())
	}
	return fmt.Sprintf("%s%s: %s", selectedItem.commitType, scopeStr, m.message.Value())
}

func (m Model) GetCommitResult() *git.CommitResult {
	return m.gitResult
}
