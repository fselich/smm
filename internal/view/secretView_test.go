package view

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SecretViewTestSuite struct {
	suite.Suite
	secretView SecretView
	width      int
	height     int
}

func (suite *SecretViewTestSuite) SetupTest() {
	suite.width = 80
	suite.height = 24
	suite.secretView = NewSecretView(suite.width, suite.height)
}

func TestSecretViewSuite(t *testing.T) {
	suite.Run(t, new(SecretViewTestSuite))
}

func (suite *SecretViewTestSuite) TestNewSecretView() {
	t := suite.T()
	sv := NewSecretView(100, 50)

	assert.False(t, sv.Hidden)
	assert.Nil(t, sv.Secret)
	assert.False(t, sv.IsFocused)
	assert.False(t, sv.IsFiltering)
	assert.Empty(t, sv.FilterValue)
	assert.False(t, sv.isFiltered)
	assert.Empty(t, sv.content)
	assert.Equal(t, 100, sv.Width())
	assert.Equal(t, 50, sv.Height())
}

func (suite *SecretViewTestSuite) TestSetContent() {
	t := suite.T()
	testContent := "This is test content\nWith multiple lines"

	suite.secretView.SetContent(testContent)

	assert.Equal(t, testContent, suite.secretView.content)
}

func (suite *SecretViewTestSuite) TestSetWidth() {
	t := suite.T()
	newWidth := 120

	suite.secretView.SetWidth(newWidth)

	assert.Equal(t, newWidth, suite.secretView.Width())
}

func (suite *SecretViewTestSuite) TestSetHeight() {
	t := suite.T()
	newHeight := 30

	suite.secretView.SetHeight(newHeight)

	assert.Equal(t, newHeight, suite.secretView.Height())
}

func (suite *SecretViewTestSuite) TestToggleFocus() {
	t := suite.T()

	// Initially not focused
	assert.False(t, suite.secretView.IsFocused)

	// Toggle to focused
	suite.secretView.ToggleFocus()
	assert.True(t, suite.secretView.IsFocused)

	// Toggle back to not focused
	suite.secretView.ToggleFocus()
	assert.False(t, suite.secretView.IsFocused)
}

func (suite *SecretViewTestSuite) TestToggleFocus_ClearsFilteringState() {
	t := suite.T()

	suite.secretView.IsFiltering = true
	suite.secretView.isFiltered = true
	suite.secretView.FilterValue = "test"

	suite.secretView.ToggleFocus()

	assert.False(t, suite.secretView.IsFiltering)
	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestSetFilteredValue_WithValue() {
	t := suite.T()
	filterValue := "search term"

	suite.secretView.SetFilteredValue(filterValue)

	assert.True(t, suite.secretView.isFiltered)
	assert.Equal(t, filterValue, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestSetFilteredValue_WithEmptyValue() {
	t := suite.T()

	suite.secretView.SetFilteredValue("")

	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestHighlightText_NoFilter() {
	t := suite.T()
	content := "This is test content"

	suite.secretView.SetContent(content)

	result := suite.secretView.HighlightText()

	assert.Equal(t, content, result)
}

func (suite *SecretViewTestSuite) TestHighlightText_WithFilter() {
	t := suite.T()
	content := "This is test content"
	filterValue := "test"

	suite.secretView.SetContent(content)
	suite.secretView.SetFilteredValue(filterValue)

	result := suite.secretView.HighlightText()

	assert.Contains(t, result, "\033[1;37;41m"+filterValue+"\033[0m")
	assert.Contains(t, result, "This is")
	assert.Contains(t, result, "content")
}

func (suite *SecretViewTestSuite) TestUpdate_StartFiltering() {
	t := suite.T()
	suite.secretView.IsFocused = true

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	assert.True(t, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringInput() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true

	// Add character to filter
	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	assert.Equal(t, "te", suite.secretView.FilterValue)
	assert.True(t, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringBackspace() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	assert.Equal(t, "tes", suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringEscape() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.False(t, suite.secretView.IsFiltering)
	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringEnter() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.False(t, suite.secretView.IsFiltering)
	assert.True(t, suite.secretView.isFiltered)
	assert.Equal(t, "test", suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteredEscape() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.isFiltered = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.False(t, suite.secretView.IsFiltering)
	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_NotFocused() {
	t := suite.T()
	suite.secretView.IsFocused = false
	initialFiltering := suite.secretView.IsFiltering

	_, _ = suite.secretView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	assert.Equal(t, initialFiltering, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestView() {
	t := suite.T()
	content := "Test content for view"

	suite.secretView.SetContent(content)

	view := suite.secretView.View()

	assert.NotEmpty(t, view)
}
