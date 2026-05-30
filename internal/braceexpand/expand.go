package braceexpand

import (
	"errors"
	"fmt"
	"strings"
)

const MaxBatchSize = 10

var (
	ErrNestedBraces  = errors.New("nested braces are not supported")
	ErrUnclosedBrace = errors.New("unclosed brace")
)

func ExpandBraces(pattern string) ([]string, error) {
	braceStart := strings.IndexByte(pattern, '{')
	if braceStart == -1 {
		return []string{pattern}, nil
	}

	braceEnd := -1
	for i := braceStart + 1; i < len(pattern); i++ {
		switch pattern[i] {
		case '{':
			return nil, ErrNestedBraces
		case '}':
			braceEnd = i
			i = len(pattern)
		}
	}
	if braceEnd == -1 {
		return nil, ErrUnclosedBrace
	}

	group := pattern[braceStart+1 : braceEnd]
	prefix := pattern[:braceStart]
	suffix := pattern[braceEnd+1:]

	elements := strings.Split(group, ",")
	valid := make([]string, 0, len(elements))
	for _, el := range elements {
		if el != "" {
			valid = append(valid, el)
		}
	}
	if len(valid) == 0 {
		return []string{}, nil
	}

	suffixExpanded, err := ExpandBraces(suffix)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(valid)*len(suffixExpanded))
	for _, el := range valid {
		for _, s := range suffixExpanded {
			result = append(result, prefix+el+s)
		}
	}

	if len(result) > MaxBatchSize {
		return nil, fmt.Errorf("brace expansion exceeds limit of %d (got %d)", MaxBatchSize, len(result))
	}

	return result, nil
}
