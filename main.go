package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	type_       string
	description string
}

func (i item) Title() string       { return i.type_ }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.type_ + " " + i.description }

type model struct {
	list    list.Model
	scope   textinput.Model
	message textinput.Model
	step    int
	err     error
}

const (
	StepTypeSelect = iota
	StepScope
	StepMessage
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170"))

	appStyle = lipgloss.NewStyle().
		Padding(1, 2)

	promptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

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
		fmt.Fprintf(w, style.Render("❯ "))
	} else {
		fmt.Fprintf(w, "  ")
	}
	fmt.Fprint(w, style.Render(fmt.Sprintf("%-10s %s", itm.type_, itm.description)))
}

func (i itemListDelegate) Height() int {
	return 1
}

func (i itemListDelegate) Spacing() int {
	return 0
}

func (i itemListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

var _ list.ItemDelegate = itemListDelegate{}

func initialModel() model {
	// Setup items
	items := []list.Item{
		item{type_: "feat", description: "A new feature"},
		item{type_: "fix", description: "A bug fix"},
		item{type_: "docs", description: "Documentation only changes"},
		item{type_: "style", description: "Changes that do not affect the meaning of the code"},
		item{type_: "refactor", description: "A code change that neither fixes a bug nor adds a feature"},
		item{type_: "perf", description: "A code change that improves performance"},
		item{type_: "test", description: "Adding missing tests or correcting existing tests"},
		item{type_: "build", description: "Changes that affect the build system or external dependencies"},
		item{type_: "ci", description: "Changes to CI configuration files and scripts"},
		item{type_: "chore", description: "Other changes that don't modify src or test files"},
	}

	// Setup list
	delegate := itemListDelegate{}
	commitList := list.New(items, delegate, 0, 0)
	commitList.Title = "Select the type of change"
	commitList.SetFilteringEnabled(true)
	commitList.SetShowHelp(true)

	// Setup scope input
	scopeInput := textinput.New()
	scopeInput.Placeholder = "scope (optional)"
	scopeInput.CharLimit = 50
	scopeInput.Width = 30

	// Setup message input
	messageInput := textinput.New()
	messageInput.Placeholder = "commit message"
	messageInput.CharLimit = 100
	messageInput.Width = 50

	return model{
		list:    commitList,
		scope:   scopeInput,
		message: messageInput,
		step:    StepTypeSelect,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
				// Create commit
				if m.message.Value() == "" {
					return m, nil
				}
				selectedItem := m.list.SelectedItem().(item)
				scopeStr := ""
				if m.scope.Value() != "" {
					scopeStr = fmt.Sprintf("(%s)", m.scope.Value())
				}
				commitMsg := fmt.Sprintf("%s%s: %s", selectedItem.type_, scopeStr, m.message.Value())

				cmd := exec.Command("git", "commit", "-m", commitMsg)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					m.err = err
					return m, nil
				}
				return m, tea.Quit
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

func (m model) View() string {
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
		s += promptStyle.Render(fmt.Sprintf("%s%s: ", selectedItem.type_, scopeStr))
		s += m.message.View()
	}

	if m.err != nil {
		s += fmt.Sprintf("\nError: %v", m.err)
	}

	return appStyle.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
