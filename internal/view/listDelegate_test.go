package view

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ListDelegateTestSuite struct {
	suite.Suite
	delegate   *ItemDelegate
	listModel  list.Model
	testSecret Secret
}

func (suite *ListDelegateTestSuite) SetupTest() {
	suite.delegate = NewListDelegate()
	suite.listModel = list.New([]list.Item{}, suite.delegate, 80, 24)
	suite.testSecret = NewSecret("test-secret", "projects/test/secrets/test-secret", "current", 1, time.Now())
}

func TestListDelegateSuite(t *testing.T) {
	suite.Run(t, new(ListDelegateTestSuite))
}

func (suite *ListDelegateTestSuite) TestNewListDelegate() {
	t := suite.T()
	delegate := NewListDelegate()

	assert.NotNil(t, delegate)
	assert.NotNil(t, delegate.Styles)
}

func (suite *ListDelegateTestSuite) TestHeight() {
	t := suite.T()

	height := suite.delegate.Height()

	assert.Equal(t, 1, height)
}

func (suite *ListDelegateTestSuite) TestSpacing() {
	t := suite.T()

	spacing := suite.delegate.Spacing()

	assert.Equal(t, 0, spacing)
}

func (suite *ListDelegateTestSuite) TestUpdate() {
	t := suite.T()
	msg := tea.KeyMsg{Type: tea.KeyDown}

	cmd := suite.delegate.Update(msg, &suite.listModel)

	assert.Nil(t, cmd)
}

func (suite *ListDelegateTestSuite) TestRender_CurrentSecret() {
	t := suite.T()
	currentSecret := NewSecret("current-secret", "path", "current", 1, time.Now())
	suite.listModel.SetItems([]list.Item{currentSecret})
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, currentSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "current-secret")
}

func (suite *ListDelegateTestSuite) TestRender_VersionSecret_FirstVersion() {
	t := suite.T()
	versionSecret := NewSecret("1", "path", "version", 1, time.Now())
	suite.listModel.SetItems([]list.Item{versionSecret})
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, versionSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "└──")
	assert.Contains(t, result, "[v.1]")
	assert.Contains(t, result, versionSecret.CreatedAt().Format("2006-01-02 15:04:05"))
}

func (suite *ListDelegateTestSuite) TestRender_VersionSecret_OtherVersion() {
	t := suite.T()
	versionSecret := NewSecret("2", "path", "version", 2, time.Now())
	suite.listModel.SetItems([]list.Item{versionSecret})
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, versionSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "├──")
	assert.Contains(t, result, "[v.2]")
	assert.Contains(t, result, versionSecret.CreatedAt().Format("2006-01-02 15:04:05"))
}

func (suite *ListDelegateTestSuite) TestRender_SelectedItem() {
	t := suite.T()
	secret := NewSecret("selected-secret", "path", "current", 1, time.Now())
	suite.listModel.SetItems([]list.Item{secret})
	suite.listModel.Select(0) // Select the first (and only) item
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, secret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "selected-secret")
}

func (suite *ListDelegateTestSuite) TestRender_UnselectedItem() {
	t := suite.T()
	secret1 := NewSecret("unselected-secret", "path", "current", 1, time.Now())
	secret2 := NewSecret("selected-secret", "path", "current", 1, time.Now())
	suite.listModel.SetItems([]list.Item{secret1, secret2})
	suite.listModel.Select(1) // Select the second item, so first is unselected
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, secret1)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "unselected-secret")
}

func (suite *ListDelegateTestSuite) TestRender_LongTitle() {
	t := suite.T()
	longTitle := "this-is-a-very-long-secret-name-that-should-be-truncated-because-it-exceeds-the-available-width"
	longSecret := NewSecret(longTitle, "path", "current", 1, time.Now())
	suite.listModel.SetItems([]list.Item{longSecret})
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, longSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	// The title should be truncated and contain ellipsis
	assert.Contains(t, result, "…")
}

func (suite *ListDelegateTestSuite) TestRender_MultipleItems() {
	t := suite.T()
	secret1 := NewSecret("secret-one", "path1", "current", 1, time.Now())
	secret2 := NewSecret("secret-two", "path2", "current", 1, time.Now())
	suite.listModel.SetItems([]list.Item{secret1, secret2})
	var output1, output2 strings.Builder

	// Render first item (selected by default)
	suite.delegate.Render(&output1, suite.listModel, 0, secret1)
	// Render second item (unselected)
	suite.delegate.Render(&output2, suite.listModel, 1, secret2)

	result1 := output1.String()
	result2 := output2.String()

	assert.NotEmpty(t, result1)
	assert.NotEmpty(t, result2)
	assert.Contains(t, result1, "secret-one")
	assert.Contains(t, result2, "secret-two")
}

func (suite *ListDelegateTestSuite) TestRender_EmptyModel() {
	t := suite.T()
	emptyModel := list.New([]list.Item{}, suite.delegate, 80, 24)
	secret := NewSecret("test-secret", "path", "current", 1, time.Now())
	var output strings.Builder

	suite.delegate.Render(&output, emptyModel, 0, secret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "test-secret")
}

func TestEllipsisConstant(t *testing.T) {
	assert.Equal(t, "…", ellipsis)
}

func (suite *ListDelegateTestSuite) TestRender_VersionWithCustomTime() {
	t := suite.T()
	customTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	versionSecret := NewSecret("3", "path", "version", 3, customTime)
	suite.listModel.SetItems([]list.Item{versionSecret})
	var output strings.Builder

	suite.delegate.Render(&output, suite.listModel, 0, versionSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "2023-12-25 15:30:45")
	assert.Contains(t, result, "[v.3]")
}

func (suite *ListDelegateTestSuite) TestRender_NarrowWidth() {
	t := suite.T()
	narrowModel := list.New([]list.Item{}, suite.delegate, 20, 24)
	longSecret := NewSecret("very-long-secret-name", "path", "current", 1, time.Now())
	narrowModel.SetItems([]list.Item{longSecret})
	var output strings.Builder

	suite.delegate.Render(&output, narrowModel, 0, longSecret)

	result := output.String()
	assert.NotEmpty(t, result)
	// Should be truncated due to narrow width
	assert.Contains(t, result, "…")
}
