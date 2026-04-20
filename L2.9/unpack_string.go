package unpack_string

// Написать функцию Go, осуществляющую примитивную распаковку строки, содержащей повторяющиеся символы/руны.

// Примеры работы функции:

//     Вход: "a4bc2d5e"
//     Выход: "aaaabccddddde"

//     Вход: "abcd"
//     Выход: "abcd" (нет цифр — ничего не меняется)

//     Вход: "45"
//     Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)

//     Вход: ""
//     Выход: "" (пустая строка -> пустая строка)

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
)

func UnpackString(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	var result strings.Builder
	result.Grow(len(str))

	var prevRune rune
	var isEscaped bool

	for _, currentRune := range str {
		if isEscaped {
			result.WriteRune(currentRune)
			prevRune = currentRune
			isEscaped = false
			continue
		}

		if currentRune == '\\' {
			isEscaped = true
			prevRune = 0
			continue
		}

		if unicode.IsDigit(currentRune) {
			if prevRune == 0 {
				return "", ErrInvalidString
			}

			count := int(currentRune - '0')

			for k := 0; k < count-1; k++ {
				result.WriteRune(prevRune)
			}
			prevRune = 0
		} else {
			result.WriteRune(currentRune)
			prevRune = currentRune
		}
	}

	if isEscaped {
		return "", ErrInvalidString
	}

	return result.String(), nil
}
