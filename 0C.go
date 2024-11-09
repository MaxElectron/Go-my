//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func parseFile(wordMultiset *map[string]int, fileName string) {
	fileBinary, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	file := strings.Split(string(fileBinary), "\n")

	for _, word := range file {
		(*wordMultiset)[word]++
	}
}

func outputRepeatedWords(wordMultiset *map[string]int) {
	for word, count := range *wordMultiset {
		if count >= 2 {
			fmt.Printf("%d\t%s\n", count, word)
		}
	}
}

func main() {
	wordMultiset := make(map[string]int)

	for _, fileName := range os.Args[1:] {
		parseFile(&wordMultiset, fileName)
	}

	outputRepeatedWords(&wordMultiset)
}
