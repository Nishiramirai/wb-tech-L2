package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Error. Number of arguments must be greater than 3")
	}

	filename := os.Args[1]
	pattern := os.Args[2]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error. Can't open file %s\n", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			fmt.Println(scanner.Text())
		}
	}
}
