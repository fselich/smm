package view

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ToastTestSuite struct {
	suite.Suite
	toast Toast
}

func (suite *ToastTestSuite) SetupTest() {
	suite.toast = NewToast()
}

func TestToastSuite(t *testing.T) {
	suite.Run(t, new(ToastTestSuite))
}

func (suite *ToastTestSuite) TestNewToast() {
	toast := NewToast()

	assert.Empty(suite.T(), toast.text)
	assert.Zero(suite.T(), toast.width)
	assert.True(suite.T(), toast.timer.IsZero())
}

func (suite *ToastTestSuite) TestSetText() {
	testText := "Test message"
	beforeTime := time.Now()

	suite.toast.SetText(testText)

	assert.Equal(suite.T(), testText, suite.toast.text)
	assert.True(suite.T(), suite.toast.timer.After(beforeTime) || suite.toast.timer.Equal(beforeTime))
}

func (suite *ToastTestSuite) TestSetWith() {
	testWidth := 100

	suite.toast.SetWith(testWidth)

	assert.Equal(suite.T(), testWidth, suite.toast.width)
}

func (suite *ToastTestSuite) TestUpdate_ShowToastMessage() {
	testText := "Update test message"
	msg := ShowToast{Text: testText}

	updatedToast, cmd := suite.toast.Update(msg)

	assert.Equal(suite.T(), testText, updatedToast.text)
	assert.Nil(suite.T(), cmd)
}

func (suite *ToastTestSuite) TestUpdate_OtherMessage() {
	suite.toast.SetText("Initial text")
	initialText := suite.toast.text

	updatedToast, cmd := suite.toast.Update("some other message")

	assert.Equal(suite.T(), initialText, updatedToast.text)
	assert.Nil(suite.T(), cmd)
}

func (suite *ToastTestSuite) TestUpdate_TextClearedAfter3Seconds() {
	suite.toast.SetText("Test message")
	suite.toast.timer = time.Now().Add(-4 * time.Second)

	updatedToast, _ := suite.toast.Update("any message")

	assert.Empty(suite.T(), updatedToast.text)
}

func (suite *ToastTestSuite) TestView_WithText() {
	testText := "View test message"
	testWidth := 50

	suite.toast.SetText(testText)
	suite.toast.SetWith(testWidth)

	view := suite.toast.View()

	assert.NotEmpty(suite.T(), view)
	assert.Contains(suite.T(), view, testText)
}

func (suite *ToastTestSuite) TestView_TextClearedAfter3Seconds() {
	suite.toast.SetText("Test message")
	suite.toast.SetWith(50)
	suite.toast.timer = time.Now().Add(-4 * time.Second)

	view := suite.toast.View()

	assert.Empty(suite.T(), suite.toast.text)
	assert.NotContains(suite.T(), view, "Test message")
}

func (suite *ToastTestSuite) TestView_EmptyText() {
	suite.toast.SetWith(50)

	view := suite.toast.View()

	assert.NotEmpty(suite.T(), view)
}

func TestToastTimerBehavior(t *testing.T) {
	toast := NewToast()
	toast.SetText("Timer test")

	view1 := toast.View()
	assert.Contains(t, view1, "Timer test")

	toast.timer = time.Now().Add(-2 * time.Second)
	view2 := toast.View()
	assert.Contains(t, view2, "Timer test")

	toast.timer = time.Now().Add(-4 * time.Second)
	view3 := toast.View()
	assert.Empty(t, toast.text)
	assert.NotContains(t, view3, "Timer test")
}
