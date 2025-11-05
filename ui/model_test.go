package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/denysvitali/git-cc/pkg/git"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	// Check that model is properly initialized
	if model.step != StepTypeSelect {
		t.Errorf("Expected StepTypeSelect (%d), got %d", StepTypeSelect, model.step)
	}

	if model.list.Title != "Select the type of change" {
		t.Errorf("Expected title 'Select the type of change', got '%s'", model.list.Title)
	}

	if !model.list.FilteringEnabled() {
		t.Error("Expected filtering to be enabled")
	}

	if model.scope.Placeholder != "scope (optional)" {
		t.Errorf("Expected placeholder 'scope (optional)', got '%s'", model.scope.Placeholder)
	}

	if model.message.Placeholder != "commit message" {
		t.Errorf("Expected placeholder 'commit message', got '%s'", model.message.Placeholder)
	}

	if model.showError != false {
		t.Error("Expected showError to be false")
	}
}

func TestModelInit(t *testing.T) {
	model := InitialModel()
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected nil command from Init")
	}
}

func TestModelUpdate_TypeSelection(t *testing.T) {
	model := InitialModel()

	// Test selecting an item and pressing Enter
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command from Update for StepScope focus")
	}

	newModelTyped := newModel.(Model)
	if newModelTyped.step != StepScope {
		t.Errorf("Expected StepScope (%d), got %d", StepScope, newModelTyped.step)
	}

	if !newModelTyped.scope.Focused() {
		t.Error("Expected scope input to be focused")
	}
}

func TestModelUpdate_ScopeInput(t *testing.T) {
	model := InitialModel()
	model.step = StepScope
	model.scope.Focus()

	// Test entering scope and pressing Enter
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command from Update for StepMessage focus")
	}

	newModelTyped := newModel.(Model)
	if newModelTyped.step != StepMessage {
		t.Errorf("Expected StepMessage (%d), got %d", StepMessage, newModelTyped.step)
	}

	if newModelTyped.scope.Focused() {
		t.Error("Expected scope input to be blurred")
	}

	if !newModelTyped.message.Focused() {
		t.Error("Expected message input to be focused")
	}
}

func TestModelUpdate_MessageInputEmpty(t *testing.T) {
	model := InitialModel()
	model.step = StepMessage
	model.message.Focus()

	// Test pressing Enter with empty message
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := model.Update(keyMsg)

	if cmd != nil {
		t.Error("Expected nil command when message is empty")
	}

	newModelTyped := newModel.(Model)
	if newModelTyped.step != StepMessage {
		t.Error("Expected to stay on StepMessage when message is empty")
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	model := InitialModel()

	// Test Ctrl+C returns a command
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected a command for Ctrl+C")
	}

	// Test 'q' key returns a command
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd = model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected a command for 'q' key")
	}
}

func TestModelUpdate_ErrorRetry(t *testing.T) {
	model := InitialModel()
	model.step = StepError
	model.showError = true

	// Test 'r' key for retry
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected command from Update for retry")
	}

	newModelTyped := newModel.(Model)
	if newModelTyped.step != StepMessage {
		t.Errorf("Expected StepMessage (%d), got %d", StepMessage, newModelTyped.step)
	}

	if newModelTyped.showError {
		t.Error("Expected showError to be false during retry setup")
	}
}

func TestBuildCommitMessage(t *testing.T) {
	model := InitialModel()
	model.step = StepMessage

	// Set up selected item
	items := []list.Item{
		item{commitType: "feat", description: "A new feature"},
	}
	model.list.SetItems(items)
	model.list.Select(0)

	// Test without scope
	model.scope.SetValue("")
	model.message.SetValue("add new feature")
	expected := "feat: add new feature"
	result := model.buildCommitMessage()
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with scope
	model.scope.SetValue("auth")
	expected = "feat(auth): add new feature"
	result = model.buildCommitMessage()
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestModelView(t *testing.T) {
	model := InitialModel()

	// Test StepTypeSelect view
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view for StepTypeSelect")
	}

	// Test StepScope view
	model.step = StepScope
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for StepScope")
	}

	// Test StepMessage view
	model.step = StepMessage
	// Set up selected item for proper display
	items := []list.Item{
		item{commitType: "feat", description: "A new feature"},
	}
	model.list.SetItems(items)
	model.list.Select(0)
	model.scope.SetValue("auth")
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for StepMessage")
	}

	// Test StepError view
	model.step = StepError
	model.showError = true
	model.gitResult = &git.CommitResult{
		Success: false,
		Message: "Test error message",
	}
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for StepError")
	}
}

func TestGetCommitResult(t *testing.T) {
	model := InitialModel()

	// Test with nil result
	if result := model.GetCommitResult(); result != nil {
		t.Error("Expected nil commit result initially")
	}

	// Test with set result
	expectedResult := &git.CommitResult{
		Success: true,
		Message: "Test success",
	}
	model.gitResult = expectedResult

	if result := model.GetCommitResult(); result != expectedResult {
		t.Error("Expected set commit result to be returned")
	}
}

func TestItemListDelegate(t *testing.T) {
	delegate := itemListDelegate{}

	// Test Height
	if delegate.Height() != 1 {
		t.Errorf("Expected height 1, got %d", delegate.Height())
	}

	// Test Spacing
	if delegate.Spacing() != 0 {
		t.Errorf("Expected spacing 0, got %d", delegate.Spacing())
	}

	// Test Update (should return nil)
	model := list.New(nil, delegate, 0, 0)
	cmd := delegate.Update(tea.KeyMsg{}, &model)
	if cmd != nil {
		t.Error("Expected nil command from delegate Update")
	}
}

func TestItemMethods(t *testing.T) {
	testItem := item{
		commitType:  "feat",
		description: "A new feature",
	}

	if testItem.Title() != "feat" {
		t.Errorf("Expected title 'feat', got '%s'", testItem.Title())
	}

	if testItem.Description() != "A new feature" {
		t.Errorf("Expected description 'A new feature', got '%s'", testItem.Description())
	}

	expectedFilter := "feat A new feature"
	if testItem.FilterValue() != expectedFilter {
		t.Errorf("Expected filter value '%s', got '%s'", expectedFilter, testItem.FilterValue())
	}
}
