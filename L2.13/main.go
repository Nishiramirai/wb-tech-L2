package main

import "flag"

// Реализовать утилиту, которая считывает входные данные (STDIN)
// и разбивает каждую строку по заданному разделителю, после чего выводит определённые поля (колонки).

type Options struct {
	fields    string // -f
	delimeter string // -d
	separated bool   // -s
}

func main() {

}

func getOpts() *Options {
	opts := &Options{}

	flag.StringVar(&opts.fields, "f", "", "usage")
	flag.StringVar(&opts.delimeter, "f", "", "usage")
	flag.BoolVar(&opts.separated, "s", false, "usage")
	flag.Parse()

	return opts
}
