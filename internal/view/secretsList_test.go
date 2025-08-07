package view

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	client "smm/internal/client"
)

type SecretsListTestSuite struct {
	suite.Suite
	fake         faker.Faker
	secretsList  SecretsList
	mockClient   client.FakeClient
	width        int
	height       int
}

func (suite *SecretsListTestSuite) SetupTest() {
	suite.fake = faker.New()
	suite.width = 80
	suite.height = 24
	suite.mockClient, _ = client.NewFakeClient("test-project")
	suite.secretsList = NewSecretsList(suite.width, suite.height, suite.mockClient)
}

func TestSecretsListSuite(t *testing.T) {
	suite.Run(t, new(SecretsListTestSuite))
}

func (suite *SecretsListTestSuite) TestNewSecret() {
	t := suite.T()
	title := "test-secret"
	fullPath := "projects/test/secrets/test-secret"
	secretType := "current"
	version := 1
	createdAt := time.Now()

	secret := NewSecret(title, fullPath, secretType, version, createdAt)

	assert.Equal(t, title, secret.title)
	assert.Equal(t, fullPath, secret.fullPath)
	assert.Equal(t, "test", secret.description)
	assert.Equal(t, secretType, secret.secretType)
	assert.Equal(t, version, secret.version)
	assert.Equal(t, createdAt, secret.createdAt)
	assert.Equal(t, 0, secret.index)
	assert.Nil(t, secret.related)
}

func (suite *SecretsListTestSuite) TestSecretFilterValue_Current() {
	t := suite.T()
	secret := NewSecret("test-secret", "path", "current", 1, time.Now())

	filterValue := secret.FilterValue()

	assert.Equal(t, "test-secret", filterValue)
}

func (suite *SecretsListTestSuite) TestSecretFilterValue_Version() {
	t := suite.T()
	related := NewSecret("related-secret", "path", "current", 1, time.Now())
	secret := NewSecret("test-secret", "path", "version", 1, time.Now())
	secret.SetRelated(&related)

	filterValue := secret.FilterValue()

	assert.Equal(t, "related-secret", filterValue)
}

func (suite *SecretsListTestSuite) TestSecretMethods() {
	t := suite.T()
	title := "test-secret"
	fullPath := "projects/test/secrets/test-secret"
	secretType := "current"
	version := 1
	createdAt := time.Now()

	secret := NewSecret(title, fullPath, secretType, version, createdAt)

	assert.Equal(t, title, secret.Title())
	assert.Equal(t, "test", secret.Description())
	assert.Equal(t, fullPath, secret.FullPath())
	assert.Equal(t, secretType, secret.Type())
	assert.Equal(t, version, secret.Version())
	assert.Equal(t, createdAt, secret.CreatedAt())
	assert.Equal(t, 0, secret.Index())
	assert.Nil(t, secret.Related())
}

func (suite *SecretsListTestSuite) TestSecretSetIndex() {
	t := suite.T()
	secret := NewSecret("test", "path", "current", 1, time.Now())
	newIndex := 5

	secret.SetIndex(newIndex)

	assert.Equal(t, newIndex, secret.Index())
}

func (suite *SecretsListTestSuite) TestSecretSetRelated() {
	t := suite.T()
	secret := NewSecret("test", "path", "current", 1, time.Now())
	related := NewSecret("related", "path", "current", 1, time.Now())

	secret.SetRelated(&related)

	assert.Equal(t, &related, secret.Related())
}

func (suite *SecretsListTestSuite) TestSecretHash() {
	t := suite.T()
	secret := NewSecret("test-secret", "path", "current", 1, time.Now())

	hash := secret.Hash()

	assert.Contains(t, hash, ".env")
	assert.Greater(t, len(hash), 4)
}

func (suite *SecretsListTestSuite) TestNewSecretsList() {
	t := suite.T()
	mockClient, _ := client.NewFakeClient("test-project")

	secretsList := NewSecretsList(100, 50, mockClient)

	assert.Equal(t, 100, secretsList.Width())
	assert.Equal(t, 50, secretsList.Height())
	assert.True(t, secretsList.IsFocused)
	assert.Empty(t, secretsList.SearchQuery)
}

func (suite *SecretsListTestSuite) TestNewSecretsList_NilClient() {
	t := suite.T()

	secretsList := NewSecretsList(100, 50, nil)

	assert.Equal(t, 100, secretsList.Width())
	assert.Equal(t, 50, secretsList.Height())
	assert.True(t, secretsList.IsFocused)
}

func (suite *SecretsListTestSuite) TestSelectedItem() {
	t := suite.T()
	// FakeClient returns predefined secrets
	secretsList := NewSecretsList(80, 24, suite.mockClient)

	selected := secretsList.SelectedItem()

	// FakeClient generates random secrets, so just check it's not empty
	assert.NotEmpty(t, selected.Title())
	assert.Equal(t, 0, selected.Index())
}

func (suite *SecretsListTestSuite) TestSelectedItem_NoItem() {
	t := suite.T()
	// Create empty client for this test
	secretsList := NewSecretsList(80, 24, nil)

	selected := secretsList.SelectedItem()

	assert.Empty(t, selected.Title())
}

func (suite *SecretsListTestSuite) TestSetWidth() {
	t := suite.T()
	newWidth := 120

	suite.secretsList.SetWidth(newWidth)

	assert.Equal(t, newWidth, suite.secretsList.Width())
}

func (suite *SecretsListTestSuite) TestSetHeight() {
	t := suite.T()
	newHeight := 30

	suite.secretsList.SetHeight(newHeight)

	assert.Equal(t, newHeight, suite.secretsList.Height())
}

func (suite *SecretsListTestSuite) TestSelect() {
	t := suite.T()
	// FakeClient returns predefined secrets
	secretsList := NewSecretsList(80, 24, suite.mockClient)
	targetIndex := 2

	secretsList.Select(targetIndex)

	assert.Equal(t, targetIndex, secretsList.Index())
}

func (suite *SecretsListTestSuite) TestSelectByName() {
	t := suite.T()
	// Create specific client for this test
	secretsList := NewSecretsList(80, 24, suite.mockClient)
	
	// Get first secret name from fake client
	firstSecret := secretsList.teaView.Items()[0].(Secret).Title()

	secretsList.SelectByName(firstSecret)

	selected := secretsList.SelectedItem()
	assert.Equal(t, firstSecret, selected.Title())
}

func (suite *SecretsListTestSuite) TestSelectByName_NotFound() {
	t := suite.T()
	// FakeClient returns predefined secrets
	secretsList := NewSecretsList(80, 24, suite.mockClient)

	secretsList.SelectByName("nonexistent-secret")

	assert.Equal(t, 0, secretsList.Index())
}

func (suite *SecretsListTestSuite) TestUpdate_NotFocused() {
	t := suite.T()
	suite.secretsList.IsFocused = false
	originalIndex := suite.secretsList.Index()

	updatedList, _ := suite.secretsList.Update(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, originalIndex, updatedList.Index())
}

func (suite *SecretsListTestSuite) TestUpdate_EscapeKey() {
	t := suite.T()
	suite.secretsList.IsFocused = true

	_, cmd := suite.secretsList.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.NotNil(t, cmd)
}

func (suite *SecretsListTestSuite) TestInsertItem() {
	secret := NewSecret("new-secret", "path", "current", 1, time.Now())

	cmd := suite.secretsList.InsertItem(0, secret)

	// InsertItem returns the command from the underlying list model
	// which might be nil or a command, so we don't assert its specific value
	_ = cmd
}

func (suite *SecretsListTestSuite) TestDelVersionItems() {
	t := suite.T()
	// Create a list with mixed item types
	items := []list.Item{
		NewSecret("current1", "path1", "current", 1, time.Now()),
		NewSecret("version1", "path2", "version", 1, time.Now()),
		NewSecret("current2", "path3", "current", 1, time.Now()),
		NewSecret("version2", "path4", "version", 2, time.Now()),
	}
	suite.secretsList.teaView.SetItems(items)

	deleted := suite.secretsList.DelVersionItems()

	assert.True(t, deleted)
	// Check that only current items remain
	remainingItems := suite.secretsList.teaView.Items()
	for _, item := range remainingItems {
		assert.Equal(t, "current", item.(Secret).Type())
	}
}

func (suite *SecretsListTestSuite) TestDelVersionItems_NoVersions() {
	t := suite.T()
	items := []list.Item{
		NewSecret("current1", "path1", "current", 1, time.Now()),
		NewSecret("current2", "path2", "current", 1, time.Now()),
	}
	suite.secretsList.teaView.SetItems(items)

	deleted := suite.secretsList.DelVersionItems()

	assert.False(t, deleted)
}

func (suite *SecretsListTestSuite) TestDeepSearch() {
	t := suite.T()
	// FakeClient will return predefined results for search

	// Test with query
	suite.secretsList.DeepSearch("query", suite.mockClient)
	assert.Equal(t, "query", suite.secretsList.SearchQuery)

	// Test with empty query
	suite.secretsList.DeepSearch("", suite.mockClient)
	assert.Empty(t, suite.secretsList.SearchQuery)
}

func (suite *SecretsListTestSuite) TestToggleFocus() {
	t := suite.T()
	initialFocus := suite.secretsList.IsFocused

	suite.secretsList.ToggleFocus()
	assert.Equal(t, !initialFocus, suite.secretsList.IsFocused)

	suite.secretsList.ToggleFocus()
	assert.Equal(t, initialFocus, suite.secretsList.IsFocused)
}

func (suite *SecretsListTestSuite) TestRealIndex() {
	t := suite.T()
	// FakeClient returns predefined secrets
	secretsList := NewSecretsList(80, 24, suite.mockClient)
	secretsList.Select(1)

	realIndex := secretsList.RealIndex()

	assert.Equal(t, 1, realIndex)
}

func (suite *SecretsListTestSuite) TestFilterMethods() {
	t := suite.T()

	isFiltering := suite.secretsList.IsFiltering()
	isFiltered := suite.secretsList.IsFiltered()
	filterValue := suite.secretsList.FilterValue()

	// These methods delegate to the underlying list
	assert.False(t, isFiltering)
	assert.False(t, isFiltered)
	assert.Empty(t, filterValue)
}

func (suite *SecretsListTestSuite) TestResetFilter() {
	// This method delegates to the underlying list
	suite.secretsList.ResetFilter()

	// No assertion needed as this is a void method that delegates
}

func (suite *SecretsListTestSuite) TestView() {
	t := suite.T()

	view := suite.secretsList.View()

	assert.NotEmpty(t, view)
}

func (suite *SecretsListTestSuite) TestInit() {
	t := suite.T()

	cmd := suite.secretsList.Init()

	assert.Nil(t, cmd)
}