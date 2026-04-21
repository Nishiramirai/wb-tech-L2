package anagram

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:  "Base case from task",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "Case sensitivity and sorting",
			input: []string{"Пятак", "тяпка", "Пятка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:     "Single words should be ignored",
			input:    []string{"привет", "мир", "стол"},
			expected: map[string][]string{},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:  "Different anagram groups",
			input: []string{"абв", "вба", "где", "едг"},
			expected: map[string][]string{
				"абв": {"абв", "вба"},
				"где": {"где", "едг"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAnagrams(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("findAnagrams() %s\ngot:  %v\nwant: %v", tt.name, got, tt.expected)
			}
		})
	}
}
