package main

import (
	"reflect"
	"testing"
)

func TestIsMatch(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		pattern string
		opts    Options
		want    bool
	}{
		{"Simple match", "hello world", "hello", Options{}, true},
		{"Case insensitive", "Hello World", "hello", Options{IgnoreCase: true}, true},
		{"Invert match", "apple", "apple", Options{Invert: true}, false},
		{"Invert no match", "banana", "apple", Options{Invert: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.Pattern = tt.pattern
			if got := isMatch(tt.line, &tt.opts); got != tt.want {
				t.Errorf("isMatch() %s: got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGrepLogic(t *testing.T) {
	lines := []string{"one", "two", "three", "four", "five"}

	t.Run("Context A1", func(t *testing.T) {
		opts := Options{Pattern: "two", After: 1}
		// Ожидаем "two", "three"
		matched := getMatchedLines(lines, &opts)
		want := []string{"two", "three"}
		if !reflect.DeepEqual(matched, want) {
			t.Errorf("got %v, want %v", matched, want)
		}
	})

	t.Run("Context B1 at start", func(t *testing.T) {
		opts := Options{Pattern: "one", Before: 1}

		matched := getMatchedLines(lines, &opts)
		want := []string{"one"}
		if !reflect.DeepEqual(matched, want) {
			t.Errorf("got %v, want %v", matched, want)
		}
	})
}

// Вспомогательная функция для тестов
func getMatchedLines(lines []string, opts *Options) []string {
	var matchedIndexes []int
	for i, line := range lines {
		if isMatch(line, opts) {
			matchedIndexes = append(matchedIndexes, i)
		}
	}

	toPrint := make(map[int]bool)
	for _, idx := range matchedIndexes {
		start := idx - opts.Before
		if start < 0 {
			start = 0
		}
		end := idx + opts.After
		if end >= len(lines) {
			end = len(lines) - 1
		}

		for i := start; i <= end; i++ {
			toPrint[i] = true
		}
	}

	var res []string
	for i := 0; i < len(lines); i++ {
		if toPrint[i] {
			res = append(res, lines[i])
		}
	}
	return res
}
