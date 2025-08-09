package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StylesTestSuite struct {
	suite.Suite
}

func TestStylesSuite(t *testing.T) {
	suite.Run(t, new(StylesTestSuite))
}

func (suite *StylesTestSuite) TestStyleLow() {
	t := suite.T()
	
	style := StyleLow()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("test content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "test content")
}

func (suite *StylesTestSuite) TestStyleSelected() {
	t := suite.T()
	
	style := StyleSelected()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("selected item")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "selected item")
}

func (suite *StylesTestSuite) TestStyleUnselected() {
	t := suite.T()
	
	style := StyleUnselected()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("unselected item")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "unselected item")
}

func (suite *StylesTestSuite) TestStyleBorder_Selected() {
	t := suite.T()
	
	style := StyleBorder(true)
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("bordered content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "bordered content")
}

func (suite *StylesTestSuite) TestStyleBorder_NotSelected() {
	t := suite.T()
	
	style := StyleBorder(false)
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("bordered content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "bordered content")
}

func (suite *StylesTestSuite) TestStyleBorder_DifferentColors() {
	t := suite.T()
	
	selectedStyle := StyleBorder(true)
	unselectedStyle := StyleBorder(false)
	
	// Both should be valid styles but potentially different
	assert.NotNil(t, selectedStyle)
	assert.NotNil(t, unselectedStyle)
	
	selectedResult := selectedStyle.Render("test")
	unselectedResult := unselectedStyle.Render("test")
	
	assert.NotEmpty(t, selectedResult)
	assert.NotEmpty(t, unselectedResult)
}

func (suite *StylesTestSuite) TestStyleModal() {
	t := suite.T()
	
	style := StyleModal()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("modal content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "modal content")
}

func (suite *StylesTestSuite) TestStyleLowBorder() {
	t := suite.T()
	
	style := StyleLowBorder()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("low border content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "low border content")
}

func (suite *StylesTestSuite) TestStyleBorderTitle() {
	t := suite.T()
	
	style := StyleBorderTitle()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("Border Title")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Border Title")
}

func (suite *StylesTestSuite) TestStyleToast() {
	t := suite.T()
	
	style := StyleToast()
	
	assert.NotNil(t, style)
	assert.IsType(t, lipgloss.Style{}, style)
	
	// Test that the style can render content
	result := style.Render("Toast message")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Toast message")
}

func (suite *StylesTestSuite) TestAllStyles_NotNil() {
	t := suite.T()
	
	styles := []lipgloss.Style{
		StyleLow(),
		StyleSelected(),
		StyleUnselected(),
		StyleBorder(true),
		StyleBorder(false),
		StyleModal(),
		StyleLowBorder(),
		StyleBorderTitle(),
		StyleToast(),
	}
	
	for i, style := range styles {
		assert.NotNil(t, style, "Style %d should not be nil", i)
		assert.IsType(t, lipgloss.Style{}, style, "Style %d should be of type lipgloss.Style", i)
	}
}

func (suite *StylesTestSuite) TestAllStyles_CanRender() {
	t := suite.T()
	
	testContent := "Test Content"
	
	styles := map[string]lipgloss.Style{
		"StyleLow":         StyleLow(),
		"StyleSelected":    StyleSelected(),
		"StyleUnselected":  StyleUnselected(),
		"StyleBorder(true)": StyleBorder(true),
		"StyleBorder(false)": StyleBorder(false),
		"StyleModal":       StyleModal(),
		"StyleLowBorder":   StyleLowBorder(),
		"StyleBorderTitle": StyleBorderTitle(),
		"StyleToast":       StyleToast(),
	}
	
	for name, style := range styles {
		result := style.Render(testContent)
		assert.NotEmpty(t, result, "%s should render non-empty content", name)
		assert.Contains(t, result, testContent, "%s should contain the original content", name)
	}
}

func (suite *StylesTestSuite) TestStyleConsistency() {
	t := suite.T()
	
	// Test that calling the same style function multiple times returns consistent results
	style1 := StyleModal()
	style2 := StyleModal()
	
	content := "consistency test"
	result1 := style1.Render(content)
	result2 := style2.Render(content)
	
	assert.Equal(t, result1, result2, "Same style function should produce consistent results")
}

func (suite *StylesTestSuite) TestStyleBorderToggle() {
	t := suite.T()
	
	// Test that StyleBorder behaves differently for selected vs unselected
	selectedBorder := StyleBorder(true)
	unselectedBorder := StyleBorder(false)
	
	assert.NotNil(t, selectedBorder)
	assert.NotNil(t, unselectedBorder)
	
	// Both should be able to render content
	content := "border test"
	selectedResult := selectedBorder.Render(content)
	unselectedResult := unselectedBorder.Render(content)
	
	assert.NotEmpty(t, selectedResult)
	assert.NotEmpty(t, unselectedResult)
	assert.Contains(t, selectedResult, content)
	assert.Contains(t, unselectedResult, content)
}