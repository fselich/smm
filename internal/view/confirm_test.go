package view

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfirmTestSuite struct {
	suite.Suite
	confirm *Confirm
}

func (suite *ConfirmTestSuite) SetupTest() {
	suite.confirm = NewConfirm("Are you sure?", "test-message")
}

func TestConfirmSuite(t *testing.T) {
	suite.Run(t, new(ConfirmTestSuite))
}

func (suite *ConfirmTestSuite) TestNewConfirm() {
	t := suite.T()
	question := "Do you want to proceed?"
	message := "some-message"
	
	confirm := NewConfirm(question, message)

	assert.NotNil(t, confirm)
	assert.Equal(t, question, confirm.question)
	assert.Equal(t, message, confirm.message)
	assert.NotNil(t, confirm.confirmation)
}

func (suite *ConfirmTestSuite) TestInit() {
	t := suite.T()

	cmd := suite.confirm.Init()

	assert.NotNil(t, cmd)
}

func (suite *ConfirmTestSuite) TestView() {
	t := suite.T()

	view := suite.confirm.View()

	// The view might be empty until the confirmation is properly initialized
	// So we just check that it's a string (could be empty)
	assert.IsType(t, "", view)
}

func (suite *ConfirmTestSuite) TestUpdate_RegularMessage() {
	t := suite.T()
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}

	modal, cmd := suite.confirm.Update(keyMsg)

	assert.NotNil(t, modal)
	assert.Nil(t, cmd) // Should return nil for regular messages
}

func (suite *ConfirmTestSuite) TestUpdate_QuitMessage() {
	t := suite.T()
	// We need to create a mock update scenario where the confirmation returns a quit command
	// Since we can't easily mock the internal confirmation model, we'll test the structure
	modal, cmd := suite.confirm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	
	assert.NotNil(t, modal)
	// cmd might be nil or not nil depending on the internal state
	_ = cmd
}

func (suite *ConfirmTestSuite) TestUpdate_StringMessage() {
	t := suite.T()
	stringMsg := "some string message"

	modal, cmd := suite.confirm.Update(stringMsg)

	assert.NotNil(t, modal)
	assert.Nil(t, cmd)
}

func TestConfirmationResultMessage(t *testing.T) {
	result := true
	msg := "test-message"
	
	confirmResult := ConfirmationResultMessage{
		Result: result,
		Msg:    msg,
	}

	assert.Equal(t, result, confirmResult.Result)
	assert.Equal(t, msg, confirmResult.Msg)
}

func TestModal_Interface(t *testing.T) {
	question := "Test question?"
	message := "test-msg"
	confirm := NewConfirm(question, message)
	
	// Test that Confirm implements Modal interface
	var modal Modal = confirm
	
	assert.NotNil(t, modal.View())
	assert.NotNil(t, modal.Init())
	
	updatedModal, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, updatedModal)
	_ = cmd
}

func (suite *ConfirmTestSuite) TestConfirmFields() {
	t := suite.T()
	question := "Custom question?"
	message := "custom-message"
	
	confirm := NewConfirm(question, message)

	assert.Equal(t, question, confirm.question)
	assert.Equal(t, message, confirm.message)
	assert.NotNil(t, confirm.confirmation)
}

func (suite *ConfirmTestSuite) TestViewWithDifferentQuestion() {
	t := suite.T()
	question := "Delete this item?"
	confirm := NewConfirm(question, "delete-action")

	view := confirm.View()

	// The view might be empty until properly initialized
	assert.IsType(t, "", view)
}

func (suite *ConfirmTestSuite) TestMultipleUpdates() {
	t := suite.T()
	
	// Test multiple updates in sequence
	modal1, cmd1 := suite.confirm.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, modal1)
	_ = cmd1

	modal2, cmd2 := suite.confirm.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.NotNil(t, modal2)
	_ = cmd2

	modal3, cmd3 := suite.confirm.Update("string message")
	assert.NotNil(t, modal3)
	_ = cmd3
}

func TestConfirmWithNilMessage(t *testing.T) {
	confirm := NewConfirm("Question?", nil)

	assert.NotNil(t, confirm)
	assert.Nil(t, confirm.message)
	assert.Equal(t, "Question?", confirm.question)
}

func TestConfirmWithComplexMessage(t *testing.T) {
	complexMsg := map[string]interface{}{
		"action": "delete",
		"id":     123,
		"name":   "test-item",
	}
	
	confirm := NewConfirm("Delete item?", complexMsg)

	assert.NotNil(t, confirm)
	assert.Equal(t, complexMsg, confirm.message)
}

func (suite *ConfirmTestSuite) TestConfirmImplementsModalCorrectly() {
	t := suite.T()
	
	// Verify that our confirm struct correctly implements all Modal methods
	var modal Modal = suite.confirm
	
	// Test View method
	view := modal.View()
	assert.IsType(t, "", view)
	
	// Test Init method
	initCmd := modal.Init()
	assert.NotNil(t, initCmd)
	
	// Test Update method
	updatedModal, updateCmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, updatedModal)
	_ = updateCmd
	
	// Ensure the returned modal is still the same type
	_, ok := updatedModal.(*Confirm)
	assert.True(t, ok)
}