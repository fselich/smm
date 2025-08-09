package view

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"smm/internal/client"
)

type SecretInfoModalTestSuite struct {
	suite.Suite
	modal      *SecretInfoModal
	secretInfo client.SecretInfo
	testSecret Secret
}

func (suite *SecretInfoModalTestSuite) SetupTest() {
	suite.secretInfo = client.SecretInfo{
		Name:       "test-secret",
		FullPath:   "projects/test-project/secrets/test-secret",
		CreateTime: time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC),
		Labels: map[string]string{
			"environment": "test",
			"team":        "engineering",
		},
		Annotations: map[string]string{
			"description": "Test secret for unit testing",
			"owner":       "test-user",
		},
	}
	suite.testSecret = NewSecret("test-secret", "projects/test-project/secrets/test-secret", "current", 1, time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC))
	suite.modal = NewSecretInfoModal(suite.secretInfo, suite.testSecret)
}

func TestSecretInfoModalSuite(t *testing.T) {
	suite.Run(t, new(SecretInfoModalTestSuite))
}

func (suite *SecretInfoModalTestSuite) TestNewSecretInfoModal() {
	t := suite.T()
	testSecret := NewSecret("test-secret", "projects/test-project/secrets/test-secret", "current", 1, time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC))
	modal := NewSecretInfoModal(suite.secretInfo, testSecret)

	assert.NotNil(t, modal)
	assert.Equal(t, suite.secretInfo, modal.secretInfo)
	assert.Equal(t, 60, modal.width)
	assert.Equal(t, 20, modal.height)
}

func (suite *SecretInfoModalTestSuite) TestInit() {
	t := suite.T()
	cmd := suite.modal.Init()
	assert.Nil(t, cmd)
}

func (suite *SecretInfoModalTestSuite) TestUpdateEscapeKey() {
	t := suite.T()
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}

	modal, cmd := suite.modal.Update(keyMsg)

	assert.NotNil(t, modal)
	assert.Nil(t, cmd)
}

func (suite *SecretInfoModalTestSuite) TestUpdateQKey() {
	t := suite.T()
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	modal, cmd := suite.modal.Update(keyMsg)

	assert.NotNil(t, modal)
	assert.Nil(t, cmd)
}

func (suite *SecretInfoModalTestSuite) TestUpdateWindowSize() {
	t := suite.T()
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 50}

	modal, cmd := suite.modal.Update(windowMsg)

	assert.NotNil(t, modal)
	assert.Nil(t, cmd)
	assert.Equal(t, 100, modal.(*SecretInfoModal).width)
	assert.Equal(t, 50, modal.(*SecretInfoModal).height)
}

func (suite *SecretInfoModalTestSuite) TestView() {
	t := suite.T()
	view := suite.modal.View()

	assert.NotEmpty(t, view)

	// Check that the view contains expected content
	assert.Contains(t, view, "Secret Information")
	assert.Contains(t, view, "test-secret")
	// The full path might be wrapped, so check for both parts
	assert.Contains(t, view, "projects/test-project/secrets/test-secr")
	assert.Contains(t, view, "et")
	assert.Contains(t, view, "2024-01-15")
	assert.Contains(t, view, "environment: test")
	assert.Contains(t, view, "team: engineering")
	assert.Contains(t, view, "description: Test secret for unit testing")
	assert.Contains(t, view, "owner: test-user")
	assert.Contains(t, view, "Press ESC to close")
}

func (suite *SecretInfoModalTestSuite) TestViewWithEmptyLabelsAndAnnotations() {
	t := suite.T()
	emptySecretInfo := client.SecretInfo{
		Name:        "empty-secret",
		FullPath:    "projects/test-project/secrets/empty-secret",
		CreateTime:  time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC),
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	}

	testSecret := NewSecret("empty-secret", "projects/test-project/secrets/empty-secret", "current", 1, time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC))
	modal := NewSecretInfoModal(emptySecretInfo, testSecret)
	view := modal.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "empty-secret")
	assert.NotContains(t, view, "Labels:")
	assert.NotContains(t, view, "Annotations:")
}

func (suite *SecretInfoModalTestSuite) TestViewAgeDisplay() {
	t := suite.T()
	now := time.Now()

	// Test recent secret (minutes ago)
	recentSecret := client.SecretInfo{
		Name:        "recent-secret",
		FullPath:    "projects/test-project/secrets/recent-secret",
		CreateTime:  now.Add(-30 * time.Minute),
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	}

	testSecretRecent := NewSecret("recent-secret", "projects/test-project/secrets/recent-secret", "current", 1, now.Add(-30*time.Minute))
	modal := NewSecretInfoModal(recentSecret, testSecretRecent)
	view := modal.View()
	assert.Contains(t, view, "minutes ago")

	// Test old secret (days ago)
	oldSecret := client.SecretInfo{
		Name:        "old-secret",
		FullPath:    "projects/test-project/secrets/old-secret",
		CreateTime:  now.Add(-5 * 24 * time.Hour),
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	}

	testSecretOld := NewSecret("old-secret", "projects/test-project/secrets/old-secret", "current", 1, now.Add(-5*24*time.Hour))
	modal = NewSecretInfoModal(oldSecret, testSecretOld)
	view = modal.View()
	assert.Contains(t, view, "5 days ago")
}

func TestSecretInfoModalImplementsModal(t *testing.T) {
	secretInfo := client.SecretInfo{
		Name:     "test",
		FullPath: "projects/test/secrets/test",
	}
	testSecret := NewSecret("test", "projects/test/secrets/test", "current", 1, time.Now())
	modal := NewSecretInfoModal(secretInfo, testSecret)

	// Ensure it implements the Modal interface
	var _ Modal = modal
}
