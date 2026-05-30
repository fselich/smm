package braceexpand

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandBraces_NoBraces(t *testing.T) {
	result, err := ExpandBraces("mi-secreto")
	require.NoError(t, err)
	assert.Equal(t, []string{"mi-secreto"}, result)
}

func TestExpandBraces_Empty(t *testing.T) {
	result, err := ExpandBraces("")
	require.NoError(t, err)
	assert.Equal(t, []string{""}, result)
}

func TestExpandBraces_Flat(t *testing.T) {
	result, err := ExpandBraces("{pro,qa,dev}-svc")
	require.NoError(t, err)
	assert.Equal(t, []string{"pro-svc", "qa-svc", "dev-svc"}, result)
}

func TestExpandBraces_Cartesian(t *testing.T) {
	result, err := ExpandBraces("{pro,qa}-{env,config}")
	require.NoError(t, err)
	assert.Equal(t, []string{"pro-env", "pro-config", "qa-env", "qa-config"}, result)
}

func TestExpandBraces_SequentialGroups(t *testing.T) {
	result, err := ExpandBraces("a{b,c}d{e,f}")
	require.NoError(t, err)
	assert.Equal(t, []string{"abde", "abdf", "acde", "acdf"}, result)
}

func TestExpandBraces_EmptyElements(t *testing.T) {
	result, err := ExpandBraces("{pro,,qa}-svc")
	require.NoError(t, err)
	assert.Equal(t, []string{"pro-svc", "qa-svc"}, result)
}

func TestExpandBraces_AllEmpty(t *testing.T) {
	result, err := ExpandBraces("{,}")
	require.NoError(t, err)
	assert.Equal(t, []string{}, result)
}

func TestExpandBraces_SingleElement(t *testing.T) {
	result, err := ExpandBraces("{a}")
	require.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
}

func TestExpandBraces_UnclosedBrace(t *testing.T) {
	_, err := ExpandBraces("{pro,qa")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnclosedBrace)
}

func TestExpandBraces_NestedBraces(t *testing.T) {
	_, err := ExpandBraces("{{a,b},c}")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNestedBraces)
}

func TestExpandBraces_ExceedsLimit(t *testing.T) {
	_, err := ExpandBraces("{a,b,c,d,e,f,g,h,i,j,k}")
	require.Error(t, err)
	assert.ErrorContains(t, err, "exceeds limit")
}

func TestExpandBraces_CartesianExceedsLimit(t *testing.T) {
	_, err := ExpandBraces("{a,b,c}*{1,2,3,4}")
	require.Error(t, err)
	assert.ErrorContains(t, err, "exceeds limit")
}

func TestExpandBraces_TripleCartesian(t *testing.T) {
	result, err := ExpandBraces("{a,b}{1,2}{x,y}")
	require.NoError(t, err)
	assert.Equal(t, []string{"a1x", "a1y", "a2x", "a2y", "b1x", "b1y", "b2x", "b2y"}, result)
}

func TestExpandBraces_NestedInSuffixAfterBrace(t *testing.T) {
	_, err := ExpandBraces("a{b,c}{{d,e}}")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNestedBraces)
}

func TestExpandBraces_PrefixBeforeBrace(t *testing.T) {
	result, err := ExpandBraces("service-{pro,qa,dev}")
	require.NoError(t, err)
	assert.Equal(t, []string{"service-pro", "service-qa", "service-dev"}, result)
}
