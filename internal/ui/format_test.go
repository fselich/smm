package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
}

func TestFormatSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}

func (suite *FormatTestSuite) TestSyntaxHighlight_PrintableContent() {
	t := suite.T()
	
	testData := []byte("KEY=value\nANOTHER_KEY=another_value")
	result := SyntaxHighlight(testData)
	
	assert.NotEmpty(t, result)
	assert.NotContains(t, result, "Non printable data")
}

func (suite *FormatTestSuite) TestSyntaxHighlight_NonPrintableContent() {
	t := suite.T()
	
	testData := []byte{0x00, 0x01, 0x02, 0x03}
	result := SyntaxHighlight(testData)
	
	assert.Contains(t, result, "Non printable data")
}

func (suite *FormatTestSuite) TestDetectFormat_EnvFormat() {
	t := suite.T()
	
	envData := []byte("KEY=value\nANOTHER_KEY=another_value")
	format := detectFormat(envData)
	
	assert.Equal(t, "bash", format)
}

func (suite *FormatTestSuite) TestDetectFormat_JSONFormat() {
	t := suite.T()
	
	jsonData := []byte(`{"key": "value", "another_key": "another_value"}`)
	format := detectFormat(jsonData)
	
	assert.Equal(t, "json", format)
}

func (suite *FormatTestSuite) TestDetectFormat_PHPFormat() {
	t := suite.T()
	
	phpData := []byte("<?php\n$config = array();")
	format := detectFormat(phpData)
	
	assert.Equal(t, "php", format)
}

func (suite *FormatTestSuite) TestDetectFormat_INIFormat() {
	t := suite.T()
	
	iniData := []byte("[section]\nkey=value\nanother_key=another_value")
	format := detectFormat(iniData)
	
	assert.Equal(t, "ini", format)
}

func (suite *FormatTestSuite) TestDetectFormat_DefaultFormat() {
	t := suite.T()
	
	plainData := []byte("This is just plain text")
	format := detectFormat(plainData)
	
	assert.Equal(t, "default", format)
}

func (suite *FormatTestSuite) TestIsEnv_ValidEnvFile() {
	t := suite.T()
	
	validEnv := "KEY=value\nANOTHER_KEY=another_value\n# This is a comment"
	result := isEnv(validEnv)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsEnv_InvalidEnvFile() {
	t := suite.T()
	
	invalidEnv := "This is not an env file"
	result := isEnv(invalidEnv)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsEnv_EmptyLines() {
	t := suite.T()
	
	envWithEmptyLines := "KEY=value\n\nANOTHER_KEY=another_value"
	result := isEnv(envWithEmptyLines)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsEnv_WithUnderscores() {
	t := suite.T()
	
	envWithUnderscores := "_KEY=value\nANOTHER_KEY_WITH_UNDERSCORES=value"
	result := isEnv(envWithUnderscores)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsJSON_ValidJSON() {
	t := suite.T()
	
	validJSON := `{"key": "value", "number": 123, "boolean": true}`
	result := isJSON(validJSON)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsJSON_InvalidJSON() {
	t := suite.T()
	
	invalidJSON := `{"key": "value", "invalid": }`
	result := isJSON(invalidJSON)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsJSON_EmptyObject() {
	t := suite.T()
	
	emptyJSON := `{}`
	result := isJSON(emptyJSON)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsPhp_ValidPHPFile() {
	t := suite.T()
	
	validPHP := "<?php\n$config = array();"
	result := isPhp(validPHP)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsPhp_InvalidPHPFile() {
	t := suite.T()
	
	invalidPHP := "$config = array();"
	result := isPhp(invalidPHP)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsPhp_WithComments() {
	t := suite.T()
	
	phpWithComments := "# Comment\n<?php\n$config = array();"
	result := isPhp(phpWithComments)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsINI_ValidINIFile() {
	t := suite.T()
	
	validINI := "[section1]\nkey1=value1\nkey2=value2\n[section2]\nkey3=value3"
	result := isINI(validINI)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsINI_InvalidINIFile() {
	t := suite.T()
	
	invalidINI := "This is not an INI file"
	result := isINI(invalidINI)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsINI_WithComments() {
	t := suite.T()
	
	iniWithComments := "# Comment\n[section]\nkey=value"
	result := isINI(iniWithComments)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsINI_EmptyLines() {
	t := suite.T()
	
	iniWithEmptyLines := "[section]\n\nkey=value\n\n"
	result := isINI(iniWithEmptyLines)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsPrintable_ValidUTF8() {
	t := suite.T()
	
	validUTF8 := []byte("Hello, world! 123 🌟")
	result := isPrintable(validUTF8)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsPrintable_InvalidUTF8() {
	t := suite.T()
	
	invalidUTF8 := []byte{0xFF, 0xFE}
	result := isPrintable(invalidUTF8)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsPrintable_ControlCharacters() {
	t := suite.T()
	
	controlChars := []byte{0x01, 0x02, 0x03}
	result := isPrintable(controlChars)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestIsPrintable_WithSpaces() {
	t := suite.T()
	
	textWithSpaces := []byte("Hello world\n\t")
	result := isPrintable(textWithSpaces)
	
	assert.True(t, result)
}

func (suite *FormatTestSuite) TestIsPrintable_DELCharacter() {
	t := suite.T()
	
	textWithDEL := []byte{72, 101, 108, 108, 111, 127} // "Hello" + DEL
	result := isPrintable(textWithDEL)
	
	assert.False(t, result)
}

func (suite *FormatTestSuite) TestSyntaxHighlight_WordWrapping() {
	t := suite.T()
	
	longLine := strings.Repeat("KEY=very_long_value_that_should_be_wrapped ", 10)
	testData := []byte(longLine)
	result := SyntaxHighlight(testData)
	
	assert.NotEmpty(t, result)
	// The result contains ANSI codes, so check for the actual content
	assert.Contains(t, result, "KEY")
	assert.Contains(t, result, "very_long_value")
}

func (suite *FormatTestSuite) TestDetectFormat_PriorityOrder() {
	t := suite.T()
	
	// Test that ENV format takes precedence over other formats when ambiguous
	ambiguousData := []byte("KEY=value")
	format := detectFormat(ambiguousData)
	
	assert.Equal(t, "bash", format, "ENV format should be detected first")
}