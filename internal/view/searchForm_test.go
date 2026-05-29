package view

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SearchFormTestSuite struct {
	suite.Suite
	searchForm *SearchForm
}

func (suite *SearchFormTestSuite) SetupTest() {
	suite.searchForm = NewSearchForm()
}

func TestSearchFormSuite(t *testing.T) {
	suite.Run(t, new(SearchFormTestSuite))
}

func (suite *SearchFormTestSuite) TestNewSearchForm() {
	t := suite.T()
	form := NewSearchForm()

	assert.NotNil(t, form)
	assert.NotNil(t, form.teaView)
	assert.Empty(t, form.Value())
}

func (suite *SearchFormTestSuite) TestValue() {
	t := suite.T()
	initialValue := suite.searchForm.Value()

	assert.Empty(t, initialValue)
}

func (suite *SearchFormTestSuite) TestInit() {
	t := suite.T()

	cmd := suite.searchForm.Init()

	assert.Nil(t, cmd)
}

func (suite *SearchFormTestSuite) TestUpdate_EnterKey() {
	t := suite.T()
	// First, simulate typing some text
	keyA := tea.KeyPressMsg{Text: "t"}
	keyE := tea.KeyPressMsg{Text: "e"}
	keyS := tea.KeyPressMsg{Text: "s"}
	keyT := tea.KeyPressMsg{Text: "t"}

	suite.searchForm.Update(keyA)
	suite.searchForm.Update(keyE)
	suite.searchForm.Update(keyS)
	suite.searchForm.Update(keyT)

	// Now test Enter key
	enterKey := tea.KeyPressMsg{Code: tea.KeyEnter}
	modal, cmd := suite.searchForm.Update(enterKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "test", suite.searchForm.Value())

	// Execute the command to get the message
	msg := cmd()
	searchMsg, ok := msg.(SearchMessage)
	assert.True(t, ok)
	assert.Equal(t, "test", searchMsg.Query)
}

func (suite *SearchFormTestSuite) TestUpdate_RegularKey() {
	t := suite.T()
	keyMsg := tea.KeyPressMsg{Text: "a"}

	modal, cmd := suite.searchForm.Update(keyMsg)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "a", suite.searchForm.Value())
}

func (suite *SearchFormTestSuite) TestUpdate_BackspaceKey() {
	t := suite.T()
	// First add some text
	keyA := tea.KeyPressMsg{Text: "a"}
	keyB := tea.KeyPressMsg{Text: "b"}
	suite.searchForm.Update(keyA)
	suite.searchForm.Update(keyB)

	// Now test backspace
	backspaceKey := tea.KeyPressMsg{Code: tea.KeyBackspace}
	modal, cmd := suite.searchForm.Update(backspaceKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Equal(t, "a", suite.searchForm.Value())
}

func (suite *SearchFormTestSuite) TestUpdate_OtherMessage() {
	t := suite.T()
	otherMsg := "some other message"

	modal, cmd := suite.searchForm.Update(otherMsg)

	assert.NotNil(t, modal)
	// cmd might be nil for non-key messages
	_ = cmd
}

func (suite *SearchFormTestSuite) TestView() {
	t := suite.T()

	view := suite.searchForm.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Query search:")
}

func (suite *SearchFormTestSuite) TestView_WithText() {
	t := suite.T()
	// Add some text first
	keyMsg := tea.KeyPressMsg{Text: "t"}
	suite.searchForm.Update(keyMsg)

	view := suite.searchForm.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Query search:")
}

func TestSearchMessage(t *testing.T) {
	query := "test query"
	msg := SearchMessage{Query: query}

	assert.Equal(t, query, msg.Query)
}

func (suite *SearchFormTestSuite) TestFormProperties() {
	t := suite.T()
	form := NewSearchForm()

	// Test that the form has the expected initial configuration
	assert.NotNil(t, form.teaView)
	assert.Empty(t, form.Value())

	// Test that we can add text and it's reflected in Value()
	keyMsg := tea.KeyPressMsg{Text: "h"}
	form.Update(keyMsg)
	keyMsg = tea.KeyPressMsg{Text: "i"}
	form.Update(keyMsg)

	assert.Equal(t, "hi", form.Value())
}

func (suite *SearchFormTestSuite) TestMultipleCharacterInput() {
	t := suite.T()
	characters := []rune{'h', 'e', 'l', 'l', 'o'}

	for _, char := range characters {
		keyMsg := tea.KeyPressMsg{Text: string(char)}
		suite.searchForm.Update(keyMsg)
	}

	assert.Equal(t, "hello", suite.searchForm.Value())
}

func (suite *SearchFormTestSuite) TestEnterWithEmptyInput() {
	t := suite.T()
	enterKey := tea.KeyPressMsg{Code: tea.KeyEnter}

	modal, cmd := suite.searchForm.Update(enterKey)

	assert.NotNil(t, modal)
	assert.NotNil(t, cmd)
	assert.Empty(t, suite.searchForm.Value())

	// Execute the command to get the message
	msg := cmd()
	searchMsg, ok := msg.(SearchMessage)
	assert.True(t, ok)
	assert.Empty(t, searchMsg.Query)
}