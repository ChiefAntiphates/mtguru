package main

import (
	"reflect"
	"testing"
)

func TestMapColors(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "All known colors",
			input:    []string{"B", "U", "G", "R", "W"},
			expected: []string{"Black", "Blue", "Green", "Red", "White"},
		},
		{
			name:     "Some unknown colors",
			input:    []string{"B", "X", "G", "Y"},
			expected: []string{"Black", "X", "Green", "Y"},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "All unknown colors",
			input:    []string{"Z", "A", "123"},
			expected: []string{"Z", "A", "123"},
		},
		{
			name:     "Mixed case sensitivity",
			input:    []string{"b", "u", "R"},
			expected: []string{"Black", "Blue", "Red"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mapColors(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("mapColors(%v) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}
