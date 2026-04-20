package unpack_string

import (
	"errors"
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:        "Basic test case 1",
			input:       "a4bc2d5e",
			expected:    "aaaabccddddde",
			expectedErr: nil,
		},
		{
			name:        "No digits",
			input:       "abcd",
			expected:    "abcd",
			expectedErr: nil,
		},
		{
			name:        "Digits only - error",
			input:       "45",
			expected:    "",
			expectedErr: ErrInvalidString,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "Single char",
			input:       "q",
			expected:    "q",
			expectedErr: nil,
		},
		{
			name:        "Single char repeated",
			input:       "a3",
			expected:    "aaa",
			expectedErr: nil,
		},
		{
			name:        "Complex string with more digits",
			input:       "qwe10rty",
			expected:    "",
			expectedErr: ErrInvalidString,
		},
		// escape - последовательности
		{
			name:        "Escaped digits",
			input:       "qwe\\4\\5",
			expected:    "qwe45",
			expectedErr: nil,
		},
		{
			name:        "Escaped and then repeated",
			input:       "qwe\\45",
			expected:    "qwe44444",
			expectedErr: nil,
		},
		{
			name:        "Escaped backslash",
			input:       "\\\\3",
			expected:    "\\\\\\",
			expectedErr: nil,
		},
		{
			name:        "Digit at start with escape - error",
			input:       "\\4",
			expected:    "4",
			expectedErr: nil,
		},
		{
			name:        "Invalid sequence - digit after escape, then digit",
			input:       "a\\3b2",
			expected:    "a3bb",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackString(tt.input)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("UnpackString(%q) error = %v, wantErr %v", tt.input, err, tt.expectedErr)
				return
			}
			if got != tt.expected {
				t.Errorf("UnpackString(%q) got = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
