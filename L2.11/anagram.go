package anagram

import (
	"sort"
	"strings"
)

func findAnagrams(words []string) map[string][]string {
	// Мапа для группировки: ключ — отсортированное слово,
	// значение — список оригинальных слов.
	groups := make(map[string][]string)

	// Мапа для хранения первого встреченного слова для каждой группы.
	// Ключ — отсортированная строка, значение — самое первое слово в нижнем регистре.
	firstWords := make(map[string]string)

	for _, word := range words {
		// Приводим к нижнему регистру
		lowerWord := strings.ToLower(word)

		// Создаем ключ (сортируем буквы)
		runes := []rune(lowerWord)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		sortedKey := string(runes)

		// Запоминаем первое встреченное слово для этой группы
		if _, exists := firstWords[sortedKey]; !exists {
			firstWords[sortedKey] = lowerWord
		}

		// Добавляем слово в группу
		groups[sortedKey] = append(groups[sortedKey], lowerWord)
	}

	// Формируем итоговый результат
	res := make(map[string][]string)

	for key, val := range groups {
		// Игнорируем группы из одного слова
		if len(val) < 2 {
			continue
		}

		// Сортируем список анаграмм по возрастанию
		sort.Strings(val)

		// Добавляем в итоговую мапу, используя первое встреченное слово как ключ
		res[firstWords[key]] = val
	}

	return res
}
