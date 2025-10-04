package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compile regexes for efficiency.
var (
	// Regex for a dice term: optional count, 'd' or 'D', sides.
	// 1: Count (optional)
	// 2: Sides
	diceRegex = regexp.MustCompile(`^(\d*)?[dD](\d+)`)

	// Regex for a modifier term: just a number.
	// 1: Value
	modifierRegex = regexp.MustCompile(`^(\d+)`)

	// Regex for an operator: + or -.
	operatorRegex = regexp.MustCompile(`^\s*([+\-])\s*`)
)

// Parser holds the state of the parsing process.
type Parser struct {
	input string // The remaining string to be parsed.
}

// NewParser creates a new parser without parsing (yet)
func NewParser(input string) *Parser {
	return &Parser{input: strings.TrimSpace(input)}
}

// Parse main entry point to parsing
func (p *Parser) Parse() (*ParsedExpression, error) {
	if p.input == "" {
		return nil, errors.New("cannot parse an empty expression")
	}

	// The first part of an expression must be a term.
	initialTerm, err := p.parseTerm()
	if err != nil {
		return nil, fmt.Errorf("invalid initial term: %w", err)
	}

	expr := &ParsedExpression{
		InitialTerm: initialTerm,
		Operations:  []Operation{},
	}

	// parse subsequent operations (+ Term, - Term).
	for p.input != "" {
		op, err := p.parseOperation()
		if err != nil {
			return nil, fmt.Errorf("failed to parse operation: %w", err)
		}
		expr.Operations = append(expr.Operations, *op)
	}

	return expr, nil
}

// parseTerm tries to parse the beginning of the input string as a Term.
func (p *Parser) parseTerm() (Term, error) {
	// First, try to match a dice roll (e.g., "d20", "2d8").
	if matches := diceRegex.FindStringSubmatch(p.input); len(matches) > 0 {
		countStr := matches[1]
		sidesStr := matches[2]

		count := 1 // Default count is 1.
		if countStr != "" {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil {
				// should be unreachable (regex validates that this should always work!)
				return nil, fmt.Errorf("invalid dice count: %s", countStr)
			}
		}

		sides, err := strconv.Atoi(sidesStr)
		if err != nil {
			// should be unreachable (regex validates that this should always work!)
			return nil, fmt.Errorf("invalid dice sides: %s", sidesStr)
		}
		if sides == 0 {
			return nil, errors.New("dice must have more than 0 sides")
		}

		// Consume matched part
		p.input = p.input[len(matches[0]):]

		return DiceRoll{Count: count, Sides: sides}, nil
	}

	// If not a dice roll, try to match a modifier (e.g., "5").
	if matches := modifierRegex.FindStringSubmatch(p.input); len(matches) > 0 {
		value, err := strconv.Atoi(matches[1])
		if err != nil {
			// Unreachable.
			return nil, fmt.Errorf("invalid modifier value: %s", matches[1])
		}

		// Consume the matched part of the string.
		p.input = p.input[len(matches[0]):]

		return Modifier{Value: value}, nil
	}

	return nil, fmt.Errorf("expected dice roll (e.g., 'd6') or modifier (e.g., '5') but found: '%s'", p.input)
}

// parseOperation parses an operator and the following term.
func (p *Parser) parseOperation() (*Operation, error) {
	matches := operatorRegex.FindStringSubmatch(p.input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("expected '+' or '-' but found: '%s'", p.input)
	}

	operator := rune(matches[1][0])

	// Consume operator + whitespace.
	p.input = p.input[len(matches[0]):]

	// term needs to follow operator!
	term, err := p.parseTerm()
	if err != nil {
		return nil, fmt.Errorf("expression ends with an operator or has invalid term after operator: %w", err)
	}

	return &Operation{Operator: operator, Term: term}, nil
}
