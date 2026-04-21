package main

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		opts     Options
		expected int // <0 если a<b, >0 если a>b, 0 если равны
	}{
		{
			name:     "Numeric Sort",
			a:        "10",
			b:        "2",
			opts:     Options{NumericSort: true},
			expected: 8, // 10 - 2
		},
		{
			name:     "Lexicographical Sort",
			a:        "10",
			b:        "2",
			opts:     Options{NumericSort: false},
			expected: -1, // "1" < "2"
		},
		{
			name:     "Column Sort (Tab separated)",
			a:        "apple\t2",
			b:        "banana\t1",
			opts:     Options{KeyColumn: 2, NumericSort: true},
			expected: 1, // 2 > 1
		},
		{
			name:     "Month Sort",
			a:        "Jan",
			b:        "Mar",
			opts:     Options{MonthSort: true},
			expected: -2, // 1 - 3
		},
		{
			name:     "Human Numeric Sort",
			a:        "1G",
			b:        "2K",
			opts:     Options{HumanNumericSort: true},
			expected: 1, // 1GB > 2KB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare(tt.a, tt.b, &tt.opts)
			// Проверяем только знак результата
			if (result < 0 && tt.expected >= 0) || (result > 0 && tt.expected <= 0) || (result == 0 && tt.expected != 0) {
				t.Errorf("compare() %s: got %v, want sign of %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestSortLines(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected []string
	}{
		{
			name:     "Reverse Numeric",
			input:    []string{"1", "10", "2"},
			opts:     Options{NumericSort: true, Reverse: true},
			expected: []string{"10", "2", "1"},
		},
		{
			name:     "Unique",
			input:    []string{"a", "a", "b", "c", "c"},
			opts:     Options{Unique: true},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := make([]string, len(tt.input))
			copy(lines, tt.input)

			sortLines(lines, &tt.opts)
			if tt.opts.Unique {
				lines = uniqueLines(lines)
			}

			if !reflect.DeepEqual(lines, tt.expected) {
				t.Errorf("%s: got %v, want %v", tt.name, lines, tt.expected)
			}
		})
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		opts     Options
		expected bool
	}{
		{
			name:     "Already Sorted",
			input:    []string{"1", "2", "10"},
			opts:     Options{NumericSort: true},
			expected: true,
		},
		{
			name:     "Disordered",
			input:    []string{"10", "2", "1"},
			opts:     Options{NumericSort: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSorted(tt.input, &tt.opts); got != tt.expected {
				t.Errorf("isSorted() %s = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
