package main

const epsilon rune = 0

func parseToNfa(l *Lexer) StateMachine {
	nfa := StateMachine{}
	nfa.Init()

ParseLoop:
	for {
		token := l.Token()
		switch token.typ {
		case LITERAL:
			nfa.addLiteral(token.content)
		case CHOOSE:
			nfa.addChoice()
		case RANGE:
			nfa.createRange(token.content)
		case LPAREN:
			nfa.startGroup()
		case RPAREN:
			nfa.endGroup()
		case REPEAT:
			nfa.addRepeat(token.content)
		case EOF:
			// EOF: think of the entire re as a group
			nfa.endGroup()
			break ParseLoop
		}
	}

	nfa.start.name = "start"
	nfa.current.accepting = true
	return nfa
}
