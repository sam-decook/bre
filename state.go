package main

import (
	"fmt"
	"strings"
)

var emptyNo int = 0

type state struct {
	name string
	accepting,
	seen bool
	transitions map[rune][]*state //to-do: change to indexing into an array
}

func (s *state) addTransition(data rune, to *state) {
	s.transitions[data] = append(s.transitions[data], to)
}

func (s state) shallow() *state {
	return &state{
		s.name,
		s.accepting,
		s.seen,
		map[rune][]*state{},
	}
}

func (s state) clone(seen map[*state]*state) *state {
	if existing := seen[&s]; existing != nil {
		return existing
	}

	cloned := s.shallow()
	// Is this correct? Or should I move it just before the return?
	// Is there be a situation where memoizing here could cause a bug?
	seen[&s] = cloned

	for ch, states := range s.transitions {
		for _, state := range states {
			next := state.clone(seen) //depth-first
			cloned.addTransition(ch, next)
		}
	}

	return cloned
}

func (s *state) print(spacing int) {
	if s.seen {
		return
	} else {
		s.seen = true
	}

	fmt.Println(strings.Repeat("-", spacing) + s.name)
	for _, states := range s.transitions {
		for _, state := range states {
			state.print(spacing + 1)
		}
	}
}

func emptyState() *state {
	description := fmt.Sprintf("empty %d", emptyNo)
	emptyNo += 1
	return &state{
		name:        description,
		accepting:   false,
		transitions: make(map[rune][]*state),
	}
}
