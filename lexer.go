package main

import (
	"fmt"
	"strings"
)

type TokenType byte

const (
	EOF TokenType = iota
	LITERAL
	CHOOSE
	RANGE
	LPAREN
	RPAREN
	REPEAT
)

// Fanciness from go's lexer: you can set the indices manually
// It's like a hashmap w/out the overhead of a hashing algorithm
var token = [...]string{
	EOF:     "EOF",
	LITERAL: "LITERAL",
	CHOOSE:  "CHOOSE",
	RANGE:   "RANGE",
	LPAREN:  "LPAREN",
	RPAREN:  "RPAREN",
	REPEAT:  "REPEAT",
}

type Token struct {
	typ     TokenType
	content string
}

func (t Token) String() string {
	return fmt.Sprintf("%7s %s", token[t.typ], t.content)
}

// Each state will lex and then return the next state to go to
// Brilliant idea from Rob Pike to replace switching through states
// https://www.youtube.com/watch?v=HxaD_trXwRE&t=2653s
// Unfortunately, I don't have many states to take advantage of it
type stateFunc func(*Lexer) stateFunc

// Initialize before use
type Lexer struct {
	regex     string //idea: use string slicing? eg: return regex[1:]
	pos       int
	ch        rune
	stateFunc stateFunc
	tokens    []Token
	tokenPos  int
}

func (l *Lexer) Init(regex string) {
	l.regex = regex
	l.pos = -1
	l.stateFunc = single
	l.tokens = make([]Token, 0, 64)
}

// Transforms the regex into a slice of tokens
func (l *Lexer) Run() {
	for l.stateFunc != nil {
		l.stateFunc = l.stateFunc(l)
	}
}

// Returns the next token
func (l *Lexer) Token() Token {
	if l.tokenPos >= len(l.tokens) {
		return Token{EOF, ""}
	}

	token := l.tokens[l.tokenPos]
	l.tokenPos += 1
	return token
}

// Returns the next token without advancing
func (l *Lexer) Peek() Token {
	if l.tokenPos >= len(l.tokens) {
		return Token{EOF, ""}
	}

	return l.tokens[l.tokenPos]
}

// Moves the lexer to the next character in the regex
func (l *Lexer) next() bool {
	if l.pos+1 >= len(l.regex) {
		return false
	}

	l.pos += 1
	l.ch = rune(l.regex[l.pos])
	return true
}

// Returns the next character in the regex without advancing
// func (l *Lexer) peek() byte {
// 	if l.pos+1 >= len(l.regex) {
// 		return 0
// 	}

// 	return l.regex[l.pos+1]
// }

func single(l *Lexer) stateFunc {
	if !l.next() {
		l.tokens = append(l.tokens, Token{EOF, ""})
		return nil
	}

	switch l.ch {
	case '|':
		l.tokens = append(l.tokens, Token{CHOOSE, "|"})
	case '(':
		l.tokens = append(l.tokens, Token{LPAREN, "("})
	case ')':
		l.tokens = append(l.tokens, Token{RPAREN, ")"})
	case '[':
		return charRange
	case '*', '+', '?', '{':
		return repeat
	default:
		l.tokens = append(l.tokens, Token{LITERAL, string(l.ch)})
	}

	return single
}

func charRange(l *Lexer) stateFunc {
	start := l.pos
	for l.next() && l.ch != ']' {
	}
	value := l.regex[start+1 : l.pos]

	l.tokens = append(l.tokens, Token{RANGE, expand(value)})
	return single
}

func expand(litRange string) string {
	set := make(map[byte]struct{})
	length := len(litRange)

	var start byte = 0x0
	for i := range length {
		if litRange[i] == '-' {
			start = litRange[i-1]
		} else if start != 0x0 {
			end := litRange[i]
			for i := start; i <= end; i++ {
				set[i] = struct{}{}
			}
			start = 0x0
		} else {
			set[litRange[i]] = struct{}{}
		}
	}

	// Collect into a string, will not be in order
	expanded := strings.Builder{}
	for k := range set {
		expanded.WriteByte(k)
	}
	return expanded.String()
}

// A token with typ = REPEAT has a value "x,y". y could be "inf"
func repeat(l *Lexer) stateFunc {
	token := Token{REPEAT, ""}
	switch l.ch {
	case '*':
		token.content = "0,"
	case '+':
		token.content = "1,"
	case '?':
		token.content = "0,1"
	case '{':
		start := l.pos
		for l.ch != '}' {
			l.next()
		}
		token.content = l.regex[start+1 : l.pos]
	}

	l.tokens = append(l.tokens, token)

	return single
}
