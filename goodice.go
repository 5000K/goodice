package goodice

import (
	"fmt"
	"math/rand"

	"github.com/5000k/goodice/parser"
)

type Goodice struct {
	expression *parser.ParsedExpression
}

func New(expression string) (*Goodice, error) {
	res, err := parser.NewParser(expression).Parse()

	if err != nil {
		return nil, err
	}

	return &Goodice{expression: res}, nil
}

type ResultType int

const (
	DiceRoll ResultType = iota
	Constant
)

type ResultPart struct {
	Type ResultType

	// Sides of a die. Will be 0 if ResultType != DiceROll
	Sides int

	// value within expression: will be negative for substractions
	Value int

	// absolute value: will always be positive
	AbsoluteValue int

	// Parts of the result (e.g. every single dice roll)
	ResultParts []int
}

type Result struct {
	Result int
	Parts  []ResultPart
}

func (g *Goodice) helpGen(term parser.Term, rand *rand.Rand) (ResultPart, error) {
	switch t := term.(type) {
	case parser.DiceRoll:
		parts := make([]int, t.Count)
		sum := 0
		for i := range parts {
			res := rand.Intn(t.Sides) + 1 // +1 since it's between 1 and n inclusively, but Intn is between 0 and n exclusively
			parts[i] = res
			sum += res
		}

		return ResultPart{
			Type:          DiceRoll,
			Sides:         t.Count,
			Value:         sum,
			AbsoluteValue: sum,
			ResultParts:   parts,
		}, nil

	case parser.Modifier:
		return ResultPart{
			Type:          Constant,
			Sides:         0,
			Value:         t.Value,
			AbsoluteValue: t.Value,
			ResultParts:   []int{t.Value},
		}, nil

	default:
		return ResultPart{}, fmt.Errorf("unknown term type: %v", term)
	}
}

// Generate will generate with a random seed
func (g *Goodice) Generate() (Result, error) {
	return g.GenerateSeeded(rand.Int())
}

// GenerateSeeded will generate with a fixed seed, for deterministic results
func (g *Goodice) GenerateSeeded(seed int) (Result, error) {
	gen := rand.New(rand.NewSource(int64(seed)))

	parts := make([]ResultPart, 1+len(g.expression.Operations))

	initial, err := g.helpGen(g.expression.InitialTerm, gen)

	if err != nil {
		return Result{}, err
	}

	currentVal := initial.Value

	parts[0] = initial

	for i, op := range g.expression.Operations {
		baseRes, err := g.helpGen(op.Term, gen)

		if err != nil {
			return Result{}, err
		}

		// apply negative
		if op.Operator == '-' {
			baseRes.Value *= -1 // negate
		}

		parts[i+1] = baseRes
		currentVal += baseRes.Value
	}

	return Result{
		Result: currentVal,
		Parts:  parts,
	}, err
}

func Generate(expression string) (Result, error) {
	return GenerateSeeded(expression, rand.Int())
}

func GenerateSeeded(expression string, seed int) (Result, error) {
	g, err := New(expression)
	if err != nil {
		return Result{}, err
	}
	return g.GenerateSeeded(seed)
}
