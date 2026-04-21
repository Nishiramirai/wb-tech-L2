package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Options содержит конфигурацию параметров сортировки.
type Options struct {
	KeyColumn            int
	NumericSort          bool
	Reverse              bool
	Unique               bool
	MonthSort            bool
	IgnoreTrailingBlanks bool
	CheckOnly            bool
	HumanNumericSort     bool
	FilePath             string
}

func main() {
	opts := getOpts()

	lines, err := readLines(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(2)
	}

	// Если установлен флаг -c, проверяем и выходим
	if opts.CheckOnly {
		if isSorted(lines, opts) {
			return
		}
		fmt.Println("sort: disorder found")
		os.Exit(1)
	}

	// Основная логика сортировки
	sortLines(lines, opts)

	// Уникальность (-u)
	if opts.Unique {
		lines = uniqueLines(lines)
	}

	// Вывод
	for _, line := range lines {
		fmt.Println(line)
	}
}

// getOpts парсит аргументы командной строки.
func getOpts() *Options {
	opts := &Options{}

	flag.IntVar(&opts.KeyColumn, "k", 0, "sort via a key (1-based column number)")
	flag.BoolVar(&opts.NumericSort, "n", false, "compare according to string numerical value")
	flag.BoolVar(&opts.Reverse, "r", false, "reverse the result of comparisons")
	flag.BoolVar(&opts.Unique, "u", false, "output only unique lines")
	flag.BoolVar(&opts.MonthSort, "M", false, "compare by month name")
	flag.BoolVar(&opts.IgnoreTrailingBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&opts.CheckOnly, "c", false, "check if sorted")
	flag.BoolVar(&opts.HumanNumericSort, "h", false, "compare human readable numbers (2K, 1G)")

	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		opts.FilePath = args[0]
	}

	return opts
}

// readLines читает данные из файла или STDIN.
func readLines(opts *Options) ([]string, error) {
	var scanner *bufio.Scanner
	if opts.FilePath != "" {
		file, err := os.Open(opts.FilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close() //nolint:errcheck
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if opts.IgnoreTrailingBlanks {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

// sortLines выполняет сортировку в зависимости от флагов.
func sortLines(lines []string, opts *Options) {
	sort.SliceStable(lines, func(i, j int) bool {
		res := compare(lines[i], lines[j], opts)
		if opts.Reverse {
			return res > 0
		}
		return res < 0
	})
}

// compare — основная логика сравнения двух строк.
func compare(a, b string, opts *Options) int {
	valA, valB := a, b

	// Извлекаем колонку, если задан -k
	if opts.KeyColumn > 0 {
		valA = getColumn(a, opts.KeyColumn)
		valB = getColumn(b, opts.KeyColumn)
	}

	// Числовая сортировка -n
	if opts.NumericSort {
		nA, _ := strconv.Atoi(valA)
		nB, _ := strconv.Atoi(valB)
		if nA != nB {
			return nA - nB
		}
	}

	// Human-readable -h
	if opts.HumanNumericSort {
		ha := parseHuman(valA)
		hb := parseHuman(valB)
		if ha != hb {
			if ha < hb {
				return -1
			}
			return 1
		}
	}

	// Месяцы -M
	if opts.MonthSort {
		mA := parseMonth(valA)
		mB := parseMonth(valB)
		if mA != mB {
			return mA - mB
		}
	}

	// Лексикографическое сравнение по умолчанию
	if valA < valB {
		return -1
	} else if valA > valB {
		return 1
	}
	return 0
}

func getColumn(line string, k int) string {
	fields := strings.Split(line, "\t")
	if k > len(fields) {
		return ""
	}
	return fields[k-1]
}

func parseMonth(s string) int {
	months := map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6,
		"JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
	}
	s = strings.ToUpper(s)
	for m, val := range months {
		if strings.HasPrefix(s, m) {
			return val
		}
	}
	return 0
}

func parseHuman(s string) float64 {
	units := map[byte]float64{'K': 1e3, 'M': 1e6, 'G': 1e9, 'T': 1e12}
	if len(s) == 0 {
		return 0
	}
	last := s[len(s)-1]
	multiplier, ok := units[last]
	if !ok {
		val, _ := strconv.ParseFloat(s, 64)
		return val
	}
	val, _ := strconv.ParseFloat(s[:len(s)-1], 64)
	return val * multiplier
}

func uniqueLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	res := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			res = append(res, lines[i])
		}
	}
	return res
}

func isSorted(lines []string, opts *Options) bool {
	for i := 1; i < len(lines); i++ {
		if opts.Reverse {
			if compare(lines[i-1], lines[i], opts) < 0 {
				return false
			}
		} else {
			if compare(lines[i-1], lines[i], opts) > 0 {
				return false
			}
		}
	}
	return true
}
