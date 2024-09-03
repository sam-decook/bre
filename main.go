package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("re2.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file.")
	}

	s := bufio.NewScanner(file)
	s.Scan()
	re := s.Text()

	lexer := Lexer{}
	lexer.Init(re)

	lexer.Run()
	fmt.Println("Lexer")
	for _, token := range lexer.tokens {
		fmt.Println(token)
	}

	nfa := parseToNfa(&lexer)
	fmt.Println("Parser")
	nfa.start.print(0)
}
