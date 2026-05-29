package view

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
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
	assert.Equal(t, -1, sv.cursorLine)
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

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "/"})

	assert.True(t, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringInput() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true

	// Add character to filter
	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "t"})
	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "e"})

	assert.Equal(t, "te", suite.secretView.FilterValue)
	assert.True(t, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringBackspace() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})

	assert.Equal(t, "tes", suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringEscape() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

	assert.False(t, suite.secretView.IsFiltering)
	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteringEnter() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.IsFiltering = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.False(t, suite.secretView.IsFiltering)
	assert.True(t, suite.secretView.isFiltered)
	assert.Equal(t, "test", suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_FilteredEscape() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.isFiltered = true
	suite.secretView.FilterValue = "test"

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

	assert.False(t, suite.secretView.IsFiltering)
	assert.False(t, suite.secretView.isFiltered)
	assert.Empty(t, suite.secretView.FilterValue)
}

func (suite *SecretViewTestSuite) TestUpdate_NotFocused() {
	t := suite.T()
	suite.secretView.IsFocused = false
	initialFiltering := suite.secretView.IsFiltering

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "/"})

	assert.Equal(t, initialFiltering, suite.secretView.IsFiltering)
}

func (suite *SecretViewTestSuite) TestView() {
	t := suite.T()
	content := "Test content for view"

	suite.secretView.SetContent(content)

	view := suite.secretView.View()

	assert.NotEmpty(t, view)
}

func (suite *SecretViewTestSuite) TestSetRawContent() {
	t := suite.T()
	raw := "DB_HOST=localhost\nDB_PORT=5432"

	suite.secretView.SetRawContent(raw)

	assert.Equal(t, raw, suite.secretView.rawContent)
}

func (suite *SecretViewTestSuite) TestEnterVisualMode() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)

	suite.secretView.EnterVisualMode()

	assert.True(t, suite.secretView.isVisualMode)
	assert.Equal(t, 0, suite.secretView.visualAnchor)
	assert.Equal(t, 0, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestEnterVisualMode_EmptyContent() {
	t := suite.T()
	suite.secretView.SetContent("")

	suite.secretView.EnterVisualMode()

	assert.False(t, suite.secretView.isVisualMode)
}

func (suite *SecretViewTestSuite) TestEnterVisualMode_AlreadyActive() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("some content")

	suite.secretView.EnterVisualMode()
	assert.True(t, suite.secretView.isVisualMode)

	suite.secretView.EnterVisualMode()
	assert.True(t, suite.secretView.isVisualMode) // still active, no re-entry
}

func (suite *SecretViewTestSuite) TestExitVisualMode() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("some content")
	suite.secretView.EnterVisualMode()
	assert.True(t, suite.secretView.isVisualMode)

	suite.secretView.ExitVisualMode()

	assert.False(t, suite.secretView.isVisualMode)
	assert.Equal(t, 0, suite.secretView.visualAnchor)
	assert.Equal(t, 0, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestExitVisualMode_NotActive() {
	t := suite.T()
	suite.secretView.ExitVisualMode() // no-op, should not panic

	assert.False(t, suite.secretView.isVisualMode)
}

func (suite *SecretViewTestSuite) TestIsVisualMode() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("content")

	assert.False(t, suite.secretView.IsVisualMode())

	suite.secretView.EnterVisualMode()

	assert.True(t, suite.secretView.IsVisualMode())
}

func (suite *SecretViewTestSuite) TestRenderContent_SingleLineSelected() {
	t := suite.T()
	content := "DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=sk-xxx"
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 1
	suite.secretView.visualCursor = 1

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Equal(t, "DB_HOST=localhost", lines[0])
	assert.Equal(t, "\033[7mDB_PORT=5432\033[27m", lines[1])
	assert.Equal(t, "API_KEY=sk-xxx", lines[2])
}

func (suite *SecretViewTestSuite) TestRenderContent_MultiLineSelected() {
	t := suite.T()
	content := "line 1\nline 2\nline 3\nline 4"
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 1
	suite.secretView.visualCursor = 3

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Equal(t, "line 1", lines[0])
	assert.Equal(t, "\033[7mline 2\033[27m", lines[1])
	assert.Equal(t, "\033[7mline 3\033[27m", lines[2])
	assert.Equal(t, "\033[7mline 4\033[27m", lines[3])
}

func (suite *SecretViewTestSuite) TestRenderContent_AnchorBelowCursor() {
	t := suite.T()
	content := "line a\nline b\nline c\nline d"
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 3
	suite.secretView.visualCursor = 1

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Equal(t, "line a", lines[0])
	assert.Equal(t, "\033[7mline b\033[27m", lines[1])
	assert.Equal(t, "\033[7mline c\033[27m", lines[2])
	assert.Equal(t, "\033[7mline d\033[27m", lines[3])
}

func (suite *SecretViewTestSuite) TestRenderContent_SyntaxHighlightedContent() {
	t := suite.T()
	content := "\033[31mBASEURL\033[0m\033[32m=qa.civitatis.com\033[0m\n\033[31mAPPLICATION_ENV\033[0m\033[32m=qa\033[0m"
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 0

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Contains(t, lines[0], "\033[7m\033[31mBASEURL")
	assert.Contains(t, lines[0], "\033[0m\033[7m\033[32m=qa.civitatis.com")
	assert.Contains(t, lines[0], "\033[7m\033[27m")
}

func (suite *SecretViewTestSuite) TestRenderContent_CursorLine() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.SetContent(content)
	suite.secretView.cursorLine = 1

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Equal(t, "line 1", lines[0])
	assert.Equal(t, "\033[4mline 2\033[24m", lines[1])
	assert.Equal(t, "line 3", lines[2])
}

func (suite *SecretViewTestSuite) TestRenderContent_CursorWithAnsiContent() {
	t := suite.T()
	content := "\033[31mRED\033[0m\n\033[32mGREEN\033[0m"
	suite.secretView.SetContent(content)
	suite.secretView.cursorLine = 1

	result := suite.secretView.renderContent()

	lines := strings.Split(result, "\n")
	assert.Contains(t, lines[1], "\033[4m\033[32mGREEN")
	assert.Contains(t, lines[1], "\033[0m\033[4m")
	assert.Contains(t, lines[1], "\033[24m")
}

func (suite *SecretViewTestSuite) TestYankSelection_SingleLine() {
	t := suite.T()
	suite.secretView.SetRawContent("DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=sk-xxx")
	suite.secretView.visualAnchor = 1
	suite.secretView.visualCursor = 1

	result := suite.secretView.yankSelection()

	assert.Equal(t, "DB_PORT=5432", result)
}

func (suite *SecretViewTestSuite) TestYankSelection_MultiLine() {
	t := suite.T()
	suite.secretView.SetRawContent("line 1\nline 2\nline 3\nline 4")
	suite.secretView.visualAnchor = 1
	suite.secretView.visualCursor = 3

	result := suite.secretView.yankSelection()

	assert.Equal(t, "line 2\nline 3\nline 4", result)
}

func (suite *SecretViewTestSuite) TestYankSelection_ReverseRange() {
	t := suite.T()
	suite.secretView.SetRawContent("a\nb\nc")
	suite.secretView.visualAnchor = 2
	suite.secretView.visualCursor = 0

	result := suite.secretView.yankSelection()

	assert.Equal(t, "a\nb\nc", result)
}

func (suite *SecretViewTestSuite) TestYankSelection_EmptyRawContent() {
	t := suite.T()
	suite.secretView.SetRawContent("")

	result := suite.secretView.yankSelection()

	assert.Empty(t, result)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeNavigateDown() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 0

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "j"})

	assert.Equal(t, 1, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeNavigateDownBoundary() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 2

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "j"})

	assert.Equal(t, 2, suite.secretView.visualCursor) // clamped
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeNavigateUp() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 2

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "k"})

	assert.Equal(t, 1, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeNavigateUpBoundary() {
	t := suite.T()
	content := "line 1\nline 2\nline 3"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 2
	suite.secretView.visualCursor = 0

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "k"})

	assert.Equal(t, 0, suite.secretView.visualCursor) // clamped
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeArrowDown() {
	t := suite.T()
	content := "line 1\nline 2"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 0

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	assert.Equal(t, 1, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeArrowUp() {
	t := suite.T()
	content := "line 1\nline 2"
	suite.secretView.IsFocused = true
	suite.secretView.SetContent(content)
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 1
	suite.secretView.visualCursor = 1

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	assert.Equal(t, 0, suite.secretView.visualCursor)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeEsc() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("some content")
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 0

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

	assert.False(t, suite.secretView.isVisualMode)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeYankExits() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("line 1\nline 2")
	suite.secretView.SetRawContent("line 1\nline 2")
	suite.secretView.isVisualMode = true
	suite.secretView.visualAnchor = 0
	suite.secretView.visualCursor = 0

	_, cmd := suite.secretView.Update(tea.KeyPressMsg{Text: "y"})

	assert.False(t, suite.secretView.isVisualMode)
	assert.NotNil(t, cmd)
}

func (suite *SecretViewTestSuite) TestUpdate_VisualModeNotFocused() {
	t := suite.T()
	suite.secretView.IsFocused = false
	suite.secretView.isVisualMode = true

	_, _ = suite.secretView.Update(tea.KeyPressMsg{Text: "j"})

	// visual mode key should not be processed when not focused
	assert.True(t, suite.secretView.isVisualMode)
}

func (suite *SecretViewTestSuite) TestToggleFocus_ExitsVisualMode() {
	t := suite.T()
	suite.secretView.IsFocused = true
	suite.secretView.SetContent("some content")
	suite.secretView.EnterVisualMode()
	assert.True(t, suite.secretView.isVisualMode)

	suite.secretView.ToggleFocus()

	assert.False(t, suite.secretView.isVisualMode)
}
