package main

import (
	"fmt"
	"strconv"
	"strings"
)

const infinity int = -100

type StateMachine struct {
	start,
	current,
	endOfChoice,
	startRepetition *state
	groups Stack[*state]
}

func (sm *StateMachine) Init() {
	sm.start = emptyState()
	sm.current = sm.start
	sm.endOfChoice = nil
	sm.startRepetition = sm.start
	sm.groups = Stack[*state]{}
	sm.groups.Init()
	sm.groups.Push(sm.start)
}

func (sm *StateMachine) addEpsilonToEmpty() {
	empty := emptyState()
	sm.current.addTransition(epsilon, empty)
	sm.current = empty
}

func (sm *StateMachine) addLiteral(content string) {
	sm.startRepetition = sm.current

	// Ensure there is only one character
	ch := rune(content[0])

	literal := emptyState()
	literal.name = "literal " + string(ch)

	sm.current.addTransition(ch, literal)
	sm.current = literal

	sm.addEpsilonToEmpty()
}

// Bug: it needs a start state with e-transitions to each choice
func (sm *StateMachine) addChoice() {
	if sm.endOfChoice == nil {
		sm.endOfChoice = sm.current
		sm.endOfChoice.name = "end of choice"
	} else {
		sm.current.addTransition(epsilon, sm.endOfChoice)
	}

	// Reset to the beginning
	sm.current = sm.groups.Peek()
	sm.current.name = "start of choice"
	sm.startRepetition = sm.current
}

func (sm *StateMachine) startGroup() {
	sm.groups.Push(sm.current)
	sm.startRepetition = sm.current
}

func (sm *StateMachine) endGroup() {
	sm.startRepetition = sm.groups.Pop()

	if sm.endOfChoice != nil {
		sm.current.addTransition(epsilon, sm.endOfChoice)
		sm.current = sm.endOfChoice
		sm.endOfChoice = nil
	}
}

func (sm *StateMachine) createRange(rangeLit string) {
	sm.groups.Push(sm.current)
	sm.startRepetition = sm.current
	sm.current.name = "start of range"

	end := emptyState()
	end.name = "end of range"
	for i := range len(rangeLit) {
		s := emptyState()
		sm.current.addTransition(rune(rangeLit[i]), s)

		s.addTransition(epsilon, end)
		s.name = fmt.Sprintf("range: %c", rangeLit[i])
	}

	sm.current = end
}

func (sm *StateMachine) addRepeat(amount string) {
	min, max := parseAmount(amount)

	if min == 0 {
		sm.startRepetition.addTransition(epsilon, sm.current)
	} else {
		// Repeat the section
		for range min - 1 {
			sm.cloneRepetition()
		}
	}

	if max == infinity {
		sm.current.addTransition(epsilon, sm.startRepetition)
	} else if max != -1 {
		// Repeat the section, and add e-transitions to the end
		end := emptyState()
		end.name = "escape"
		sm.current.addTransition(epsilon, end)
		for range max - min {
			sm.cloneRepetition()
		}
		sm.current = end
	}
}

// Extract the amount of repeats from `min[,[max]]`
func parseAmount(amount string) (min, max int) {
	parts := strings.Split(amount, ",")

	min, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("Invalid minimum number for repeat: " + parts[0])
	}

	if len(parts) == 1 {
		return min, -1
	}

	max, err = strconv.Atoi(parts[1])
	if parts[1] == "" {
		max = infinity
	} else if err != nil {
		panic("Invalid maximum number for repeat: " + parts[1])
	}

	return min, max
}

func (sm *StateMachine) cloneRepetition() {
	seen := make(map[*state]*state)

	// Recursively clone all proceeding states
	start := sm.startRepetition.clone(seen)

	// Walk graph to find the end
	end := start
	for len(end.transitions) != 0 {
		for _, states := range end.transitions {
			end = states[0]
			break //sufficient to follow one path
		}
	}

	// Add on the cloned section
	sm.current.addTransition(epsilon, start)
	sm.startRepetition = start
	sm.current = end
}
