package main

// Реализовать утилиту фильтрации текстового потока (аналог команды grep).
// Программа должна читать входной поток (STDIN или файл) и выводить строки,
// соответствующие заданному шаблону (подстроке или регулярному выражению).

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Options struct {
	After      int  // -A
	Before     int  // -B
	Context    int  // -C
	Count      bool // -c
	IgnoreCase bool // -i
	Invert     bool // -v
	Fixed      bool // -F
	LineNum    bool // -n
	Pattern    string
	FilePath   string
}

func main() {
	opts := getOpts()

	lines, err := readLines(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(2)
	}

	runGrep(lines, opts)
}

func getOpts() *Options {
	opts := &Options{}

	flag.IntVar(&opts.After, "A", 0, "show N lines after match")
	flag.IntVar(&opts.Before, "B", 0, "show N lines before match")
	flag.IntVar(&opts.Context, "C", 0, "show N lines around match")
	flag.BoolVar(&opts.Count, "c", false, "show count of matching lines")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&opts.Invert, "v", false, "invert match")
	flag.BoolVar(&opts.Fixed, "F", false, "fixed string pattern")
	flag.BoolVar(&opts.LineNum, "n", false, "show line numbers")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: grep [options] pattern [file]")
		os.Exit(2)
	}

	opts.Pattern = args[0]
	if len(args) > 1 {
		opts.FilePath = args[1]
	}

	if opts.Context > 0 {
		opts.After = opts.Context
		opts.Before = opts.Context
	}

	return opts
}

func readLines(opts *Options) ([]string, error) {
	var scanner *bufio.Scanner
	if opts.FilePath != "" {
		file, err := os.Open(opts.FilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close() // nolint:errcheck
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func runGrep(lines []string, opts *Options) {
	var matchedIndexes []int

	// Находим все индексы строк, которые подходят под шаблон
	for i, line := range lines {
		if isMatch(line, opts) {
			matchedIndexes = append(matchedIndexes, i)
		}
	}

	// Если нужен только подсчет (-c)
	if opts.Count {
		fmt.Println(len(matchedIndexes))
		return
	}

	// Вычисляем, какие строки должны быть выведены (с учетом контекста)
	// Используем map, чтобы избежать дубликатов при пересечении контекста
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

	// Вывод результата
	for i := 0; i < len(lines); i++ {
		if toPrint[i] {
			prefix := ""
			if opts.LineNum {
				prefix = fmt.Sprintf("%d:", i+1)
			}
			fmt.Printf("%s%s\n", prefix, lines[i])
		}
	}
}

func isMatch(line string, opts *Options) bool {
	match := false
	pattern := opts.Pattern
	target := line

	if opts.IgnoreCase {
		pattern = strings.ToLower(pattern)
		target = strings.ToLower(target)
	}

	if opts.Fixed {
		match = strings.Contains(target, pattern)
	} else {
		finalPattern := opts.Pattern
		if opts.IgnoreCase {
			finalPattern = "(?i)" + finalPattern
		}
		re, err := regexp.Compile(finalPattern)
		if err != nil {
			return false
		}
		match = re.MatchString(line)
	}

	if opts.Invert {
		return !match
	}
	return match
}
