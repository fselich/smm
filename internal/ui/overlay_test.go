package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OverlayTestSuite struct {
	suite.Suite
}

func TestOverlaySuite(t *testing.T) {
	suite.Run(t, new(OverlayTestSuite))
}

func (suite *OverlayTestSuite) TestGetLines_SingleLine() {
	t := suite.T()
	
	text := "Hello world"
	lines, widest := getLines(text)
	
	assert.Equal(t, []string{"Hello world"}, lines)
	assert.Equal(t, 11, widest)
}

func (suite *OverlayTestSuite) TestGetLines_MultipleLines() {
	t := suite.T()
	
	text := "Hello\nworld\ntest"
	lines, widest := getLines(text)
	
	assert.Equal(t, []string{"Hello", "world", "test"}, lines)
	assert.Equal(t, 5, widest)
}

func (suite *OverlayTestSuite) TestGetLines_DifferentWidths() {
	t := suite.T()
	
	text := "Short\nThis is a longer line\nMed"
	lines, widest := getLines(text)
	
	assert.Equal(t, []string{"Short", "This is a longer line", "Med"}, lines)
	assert.Equal(t, 21, widest)
}

func (suite *OverlayTestSuite) TestGetLines_EmptyString() {
	t := suite.T()
	
	text := ""
	lines, widest := getLines(text)
	
	assert.Equal(t, []string{""}, lines)
	assert.Equal(t, 0, widest)
}

func (suite *OverlayTestSuite) TestModalOverlay_BasicFunctionality() {
	t := suite.T()
	
	fg := "Modal content"
	bg := "Background content\nSecond line\nThird line"
	
	result := ModalOverlay(fg, bg)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Modal content")
}

func (suite *OverlayTestSuite) TestModalOverlay_ForegroundLargerThanBackground() {
	t := suite.T()
	
	fg := "Very long modal content that exceeds background\nSecond line of modal\nThird line\nFourth line"
	bg := "Short bg"
	
	result := ModalOverlay(fg, bg)
	
	assert.NotEmpty(t, result)
	// When fg is larger, it should return the styled fg
	assert.Contains(t, result, "Very long modal content")
}

func (suite *OverlayTestSuite) TestPlaceOverlay_BasicPlacement() {
	t := suite.T()
	
	fg := "FG"
	bg := "BACKGROUND\nSECONDLINE"
	
	result := PlaceOverlay(2, 1, fg, bg, false)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "FG")
}

func (suite *OverlayTestSuite) TestPlaceOverlay_WithShadow() {
	t := suite.T()
	
	fg := "Modal"
	bg := "Background text\nSecond line\nThird line\nFourth line"
	
	result := PlaceOverlay(2, 1, fg, bg, true)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Modal")
}

func (suite *OverlayTestSuite) TestPlaceOverlay_ZeroCoordinates() {
	t := suite.T()
	
	fg := "FG"
	bg := "BACKGROUND\nSECONDLINE"
	
	result := PlaceOverlay(0, 0, fg, bg, false)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "FG")
}

func (suite *OverlayTestSuite) TestPlaceOverlay_ForegroundLargerThanBackground() {
	t := suite.T()
	
	fg := "Very long foreground content\nSecond line"
	bg := "Short"
	
	result := PlaceOverlay(0, 0, fg, bg, false)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Very long foreground")
}

func (suite *OverlayTestSuite) TestCutLeft_BasicFunctionality() {
	t := suite.T()
	
	text := "Hello world"
	result := cutLeft(text, 6)
	
	assert.Equal(t, "world", result)
}

func (suite *OverlayTestSuite) TestCutLeft_WithANSICodes() {
	t := suite.T()
	
	text := "\x1b[31mHello\x1b[0m world"
	result := cutLeft(text, 6)
	
	assert.Contains(t, result, "world")
}

func (suite *OverlayTestSuite) TestCutLeft_ZeroCutWidth() {
	t := suite.T()
	
	text := "Hello world"
	result := cutLeft(text, 0)
	
	assert.Equal(t, "Hello world", result)
}

func (suite *OverlayTestSuite) TestCutLeft_CutWidthExceedsString() {
	t := suite.T()
	
	text := "Hello"
	result := cutLeft(text, 10)
	
	assert.Equal(t, "", result)
}

func (suite *OverlayTestSuite) TestClamp_WithinBounds() {
	t := suite.T()
	
	result := clamp(5, 1, 10)
	assert.Equal(t, 5, result)
}

func (suite *OverlayTestSuite) TestClamp_BelowLowerBound() {
	t := suite.T()
	
	result := clamp(-5, 1, 10)
	assert.Equal(t, 1, result)
}

func (suite *OverlayTestSuite) TestClamp_AboveUpperBound() {
	t := suite.T()
	
	result := clamp(15, 1, 10)
	assert.Equal(t, 10, result)
}

func (suite *OverlayTestSuite) TestMax_FirstLarger() {
	t := suite.T()
	
	result := max(10, 5)
	assert.Equal(t, 10, result)
}

func (suite *OverlayTestSuite) TestMax_SecondLarger() {
	t := suite.T()
	
	result := max(5, 10)
	assert.Equal(t, 10, result)
}

func (suite *OverlayTestSuite) TestMax_Equal() {
	t := suite.T()
	
	result := max(5, 5)
	assert.Equal(t, 5, result)
}

func (suite *OverlayTestSuite) TestMin_FirstSmaller() {
	t := suite.T()
	
	result := min(5, 10)
	assert.Equal(t, 5, result)
}

func (suite *OverlayTestSuite) TestMin_SecondSmaller() {
	t := suite.T()
	
	result := min(10, 5)
	assert.Equal(t, 5, result)
}

func (suite *OverlayTestSuite) TestMin_Equal() {
	t := suite.T()
	
	result := min(5, 5)
	assert.Equal(t, 5, result)
}

func (suite *OverlayTestSuite) TestWhitespace_Render_DefaultChars() {
	t := suite.T()
	
	ws := whitespace{}
	result := ws.render(5)
	
	assert.Equal(t, 5, len(result))
	assert.Equal(t, strings.Repeat(" ", 5), result)
}

func (suite *OverlayTestSuite) TestWhitespace_Render_CustomChars() {
	t := suite.T()
	
	ws := whitespace{chars: "-"}
	result := ws.render(3)
	
	assert.Equal(t, "---", result)
}

func (suite *OverlayTestSuite) TestWhitespace_Render_ZeroWidth() {
	t := suite.T()
	
	ws := whitespace{}
	result := ws.render(0)
	
	assert.Equal(t, "", result)
}

func (suite *OverlayTestSuite) TestWhitespace_Render_MultipleChars() {
	t := suite.T()
	
	ws := whitespace{chars: "ab"}
	result := ws.render(5)
	
	assert.Contains(t, result, "a")
	assert.Contains(t, result, "b")
	cleanResult := strings.ReplaceAll(result, "\x1b", "")
	assert.True(t, len(cleanResult) >= 5) // Account for potential ANSI codes
}

func (suite *OverlayTestSuite) TestPlaceOverlay_CoordinateClamping() {
	t := suite.T()
	
	fg := "FG"
	bg := "BACKGROUND"
	
	// Test negative coordinates get clamped to 0
	result := PlaceOverlay(-5, -5, fg, bg, false)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "FG")
}

func (suite *OverlayTestSuite) TestPlaceOverlay_ExcessiveCoordinates() {
	t := suite.T()
	
	fg := "X"
	bg := "BACKGROUND"
	
	// Test coordinates larger than background get clamped
	result := PlaceOverlay(100, 100, fg, bg, false)
	
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "X")
}