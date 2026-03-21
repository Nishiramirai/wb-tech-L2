package unpack_string

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("Invalid string")
)

func UnpackString(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runeStr := []rune(str)
	var prevRune rune
	isEscaped := false

	var result strings.Builder
	result.Grow(len(str))

	for i := 0; i < len(runeStr); i++ {
		currentRune := runeStr[i]

		if isEscaped {
			result.WriteRune(currentRune)
			prevRune = currentRune
			isEscaped = false
			continue
		}

		if currentRune == '\\' {
			isEscaped = true
			prevRune = rune(0)
			continue
		}

		if unicode.IsDigit(currentRune) {
			if prevRune == rune(0) {
				return "", ErrInvalidString
			}
			// Руна уже проверена на isDigit, поэтому
			// ошибку можно не обрабатывать
			count, _ := strconv.Atoi(string(currentRune))
			for k := 0; k < count-1; k++ {
				result.WriteRune(prevRune)
			}
			prevRune = rune(0)
		} else {
			result.WriteRune(currentRune)
			prevRune = currentRune
		}

	}

	return result.String(), nil
}
