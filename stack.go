package main

type Stack[E any] struct {
	data []E
}

func (s *Stack[E]) Init() Stack[E] {
	return Stack[E]{
		data: make([]E, 0, 16),
	}
}

func (s *Stack[E]) Len() int {
	return len(s.data)
}

func (s *Stack[E]) Push(element E) {
	s.data = append(s.data, element)
}

// Warning: doesn't check for empty stack
func (s *Stack[E]) Peek() E {
	return s.data[s.Len()-1]
}

// Warning: doesn't check for empty stack
func (s *Stack[E]) Pop() E {
	element := s.Peek()
	s.data = s.data[:s.Len()-1]
	return element
}
