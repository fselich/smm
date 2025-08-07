package view

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProjectSelectorTestSuite struct {
	suite.Suite
	projectSelector *ProjectSelector
}

func (suite *ProjectSelectorTestSuite) SetupTest() {
	// Reset viper for each test
	viper.Reset()
	viper.Set("projectids", []string{"project1", "project2"})
	suite.projectSelector = NewProjectSelectorModal()
}

func (suite *ProjectSelectorTestSuite) TearDownTest() {
	viper.Reset()
}

func TestProjectSelectorSuite(t *testing.T) {
	suite.Run(t, new(ProjectSelectorTestSuite))
}

func (suite *ProjectSelectorTestSuite) TestNewProjectSelectorModal() {
	t := suite.T()
	viper.Set("projectids", []string{"test1", "test2"})
	
	selector := NewProjectSelectorModal()

	assert.NotNil(t, selector)
	assert.NotNil(t, selector.teaView)
	assert.Empty(t, selector.alertText)
	assert.NotNil(t, selector.alertStyle)
	assert.Empty(t, selector.Value())
}

func (suite *ProjectSelectorTestSuite) TestInit() {
	t := suite.T()

	cmd := suite.projectSelector.Init()

	assert.Nil(t, cmd)
}

func (suite *ProjectSelectorTestSuite) TestValue() {
	t := suite.T()
	initialValue := suite.projectSelector.Value()

	assert.Empty(t, initialValue)
}

func (suite *ProjectSelectorTestSuite) TestSetAlert() {
	t := suite.T()
	alertText := "Error: Invalid project ID"

	suite.projectSelector.SetAlert(alertText)

	assert.Equal(t, alertText, suite.projectSelector.alertText)
}

func (suite *ProjectSelectorTestSuite) TestClearAlert() {
	t := suite.T()
	suite.projectSelector.SetAlert("Some alert")
	
	suite.projectSelector.ClearAlert()

	assert.Empty(t, suite.projectSelector.alertText)
}

func (suite *ProjectSelectorTestSuite) TestUpdate_EnterKey() {
	t := suite.T()
	viper.Set("projectIds", []string{"existing1", "existing2"})
	
	// First add some text
	keyT := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	keyE := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	keyS := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	keyT2 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	
	suite.projectSelector.Update(keyT)
	suite.projectSelector.Update(keyE)
	suite.projectSelector.Update(keyS)
	suite.projectSelector.Update(keyT2)

	// Now test Enter key
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	modal, cmd := suite.projectSelector.Update(enterKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "test", suite.projectSelector.Value())

	// Check that viper settings were updated
	assert.Equal(t, "test", viper.GetString("selected"))
	projectIds := viper.GetStringSlice("projectIds")
	assert.Contains(t, projectIds, "test")
	assert.Contains(t, projectIds, "existing1")
	assert.Contains(t, projectIds, "existing2")

	// Execute the command to get the message
	msg := cmd()
	projectMsg, ok := msg.(ProjectSelectedMessage)
	assert.True(t, ok)
	assert.Equal(t, "test", projectMsg.ProjectId)
}

func (suite *ProjectSelectorTestSuite) TestUpdate_EnterKey_DuplicateProject() {
	t := suite.T()
	viper.Set("projectIds", []string{"test", "other"})
	
	// Add text that matches existing project
	keyT := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	keyE := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	keyS := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	keyT2 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	
	suite.projectSelector.Update(keyT)
	suite.projectSelector.Update(keyE)
	suite.projectSelector.Update(keyS)
	suite.projectSelector.Update(keyT2)

	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	modal, cmd := suite.projectSelector.Update(enterKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)

	// Check that no duplicates were added
	projectIds := viper.GetStringSlice("projectIds")
	testCount := 0
	for _, id := range projectIds {
		if id == "test" {
			testCount++
		}
	}
	assert.Equal(t, 1, testCount) // Should only appear once
}

func (suite *ProjectSelectorTestSuite) TestUpdate_RegularKey() {
	t := suite.T()
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

	modal, cmd := suite.projectSelector.Update(keyMsg)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "a", suite.projectSelector.Value())
}

func (suite *ProjectSelectorTestSuite) TestUpdate_BackspaceKey() {
	t := suite.T()
	// First add some text
	keyA := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	keyB := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
	suite.projectSelector.Update(keyA)
	suite.projectSelector.Update(keyB)

	// Now test backspace
	backspaceKey := tea.KeyMsg{Type: tea.KeyBackspace}
	modal, cmd := suite.projectSelector.Update(backspaceKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "a", suite.projectSelector.Value())
}

func (suite *ProjectSelectorTestSuite) TestUpdate_OtherMessage() {
	t := suite.T()
	otherMsg := "some other message"

	modal, cmd := suite.projectSelector.Update(otherMsg)

	assert.NotNil(t, modal)
	// cmd might be nil for non-key messages
	_ = cmd
}

func (suite *ProjectSelectorTestSuite) TestView_WithoutAlert() {
	t := suite.T()

	view := suite.projectSelector.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Project ID:")
}

func (suite *ProjectSelectorTestSuite) TestView_WithAlert() {
	t := suite.T()
	alertText := "Invalid project!"
	suite.projectSelector.SetAlert(alertText)

	view := suite.projectSelector.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Project ID:")
	assert.Contains(t, view, alertText)
}

func (suite *ProjectSelectorTestSuite) TestView_AlertClearing() {
	t := suite.T()
	suite.projectSelector.SetAlert("Alert text")
	viewWithAlert := suite.projectSelector.View()
	assert.Contains(t, viewWithAlert, "Alert text")

	suite.projectSelector.ClearAlert()
	viewWithoutAlert := suite.projectSelector.View()
	assert.NotContains(t, viewWithoutAlert, "Alert text")
}

func TestProjectSelectedMessage(t *testing.T) {
	projectId := "my-project-123"
	msg := ProjectSelectedMessage{ProjectId: projectId}

	assert.Equal(t, projectId, msg.ProjectId)
}

func TestShowProjectSelectMsg(t *testing.T) {
	textAlert := "Please select a project"
	msg := ShowProjectSelectMsg{TextAlert: textAlert}

	assert.Equal(t, textAlert, msg.TextAlert)
}

func (suite *ProjectSelectorTestSuite) TestProjectSelectorImplementsModal() {
	t := suite.T()
	
	// Test that ProjectSelector implements Modal interface
	var modal Modal = suite.projectSelector
	
	assert.NotNil(t, modal.View())
	assert.Nil(t, modal.Init())
	
	updatedModal, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, updatedModal)
	_ = cmd
}

func (suite *ProjectSelectorTestSuite) TestMultipleCharacterInput() {
	t := suite.T()
	characters := []rune{'m', 'y', '-', 'p', 'r', 'o', 'j', 'e', 'c', 't'}

	for _, char := range characters {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		suite.projectSelector.Update(keyMsg)
	}

	assert.Equal(t, "my-project", suite.projectSelector.Value())
}

func (suite *ProjectSelectorTestSuite) TestEnterWithEmptyInput() {
	t := suite.T()
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}

	modal, cmd := suite.projectSelector.Update(enterKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Empty(t, suite.projectSelector.Value())

	// Execute the command to get the message
	msg := cmd()
	projectMsg, ok := msg.(ProjectSelectedMessage)
	assert.True(t, ok)
	assert.Empty(t, projectMsg.ProjectId)
}

func (suite *ProjectSelectorTestSuite) TestViperIntegration() {
	t := suite.T()
	initialProjects := []string{"proj1", "proj2"}
	viper.Set("projectIds", initialProjects)
	
	// Add text and press enter
	keyP := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	keyE := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	keyW := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}}
	
	suite.projectSelector.Update(keyP)
	suite.projectSelector.Update(keyE)
	suite.projectSelector.Update(keyW)
	
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	suite.projectSelector.Update(enterKey)

	// Check that viper was updated correctly
	assert.Equal(t, "new", viper.GetString("selected"))
	projectIds := viper.GetStringSlice("projectIds")
	assert.Contains(t, projectIds, "new")
	assert.Contains(t, projectIds, "proj1")
	assert.Contains(t, projectIds, "proj2")
}

func (suite *ProjectSelectorTestSuite) TestAlertStyling() {
	t := suite.T()
	
	// Test that alert style is properly initialized
	assert.NotNil(t, suite.projectSelector.alertStyle)
	
	// Set an alert and verify view changes
	suite.projectSelector.SetAlert("Test alert")
	viewWithAlert := suite.projectSelector.View()
	
	suite.projectSelector.ClearAlert()
	viewWithoutAlert := suite.projectSelector.View()
	
	// Views should be different
	assert.NotEqual(t, viewWithAlert, viewWithoutAlert)
}