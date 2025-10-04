package parser

import (
	"fmt"
	"strconv"
)

type TermType int

const (
	DiceRollType TermType = iota
	ModifierType
)

type Term interface {
	Type() TermType
	String() string
}

type DiceRoll struct {
	Count int // The 'X' part. Defaults to 1 if not present.
	Sides int // The 'Y' part.
}

func (dr DiceRoll) Type() TermType {
	return DiceRollType
}

// String creates a canonical string representation of a DiceRoll, e.g., "3D6".
func (dr DiceRoll) String() string {
	return fmt.Sprintf("%dD%d", dr.Count, dr.Sides)
}

type Modifier struct {
	Value int
}

func (m Modifier) Type() TermType {
	return ModifierType
}

// String creates a string representation of a Modifier.
func (m Modifier) String() string {
	return strconv.Itoa(m.Value)
}

type Operation struct {
	Operator rune // '+' or '-'
	Term     Term
}

type ParsedExpression struct {
	// The first term in the expression (e.g., "1D20" in "1D20+5").
	InitialTerm Term
	// A list of subsequent operations (e.g., "+5").
	Operations []Operation
}
