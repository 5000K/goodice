package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/5000K/goodice/parser"
)

func TestParser_Parse(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name               string
		input              string
		expectedExpression *parser.ParsedExpression
		expectError        bool
		errorContains      string
	}{
		// --- Valid Cases ---
		{
			name:  "simple d20",
			input: "d20",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 20},
				Operations:  []parser.Operation{},
			},
			expectError: false,
		},
		{
			name:  "simple uppercase D6",
			input: "D6",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 6},
				Operations:  []parser.Operation{},
			},
			expectError: false,
		},
		{
			name:  "multiple dice 3d6",
			input: "3d6",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 3, Sides: 6},
				Operations:  []parser.Operation{},
			},
			expectError: false,
		},
		{
			name:  "simple modifier",
			input: "5",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.Modifier{Value: 5},
				Operations:  []parser.Operation{},
			},
			expectError: false,
		},
		{
			name:  "dice plus modifier",
			input: "d20+5",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 20},
				Operations:  []parser.Operation{{Operator: '+', Term: parser.Modifier{Value: 5}}},
			},
			expectError: false,
		},
		{
			name:  "dice minus modifier with spaces",
			input: " D8 - 2 ",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 8},
				Operations:  []parser.Operation{{Operator: '-', Term: parser.Modifier{Value: 2}}},
			},
			expectError: false,
		},
		{
			name:  "modifier plus dice",
			input: "7+2d4",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.Modifier{Value: 7},
				Operations:  []parser.Operation{{Operator: '+', Term: parser.DiceRoll{Count: 2, Sides: 4}}},
			},
			expectError: false,
		},
		{
			name:  "complex expression",
			input: "d20 + 2d8 - 5",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 20},
				Operations: []parser.Operation{
					{Operator: '+', Term: parser.DiceRoll{Count: 2, Sides: 8}},
					{Operator: '-', Term: parser.Modifier{Value: 5}},
				},
			},
			expectError: false,
		},
		{
			name:  "long chain of operations",
			input: "1d100+1d20-1d12+1d4-1",
			expectedExpression: &parser.ParsedExpression{
				InitialTerm: parser.DiceRoll{Count: 1, Sides: 100},
				Operations: []parser.Operation{
					{Operator: '+', Term: parser.DiceRoll{Count: 1, Sides: 20}},
					{Operator: '-', Term: parser.DiceRoll{Count: 1, Sides: 12}},
					{Operator: '+', Term: parser.DiceRoll{Count: 1, Sides: 4}},
					{Operator: '-', Term: parser.Modifier{Value: 1}},
				},
			},
			expectError: false,
		},

		// --- Error Cases ---
		{
			name:          "empty input",
			input:         "",
			expectError:   true,
			errorContains: "empty expression",
		},
		{
			name:          "whitespace input",
			input:         "   ",
			expectError:   true,
			errorContains: "empty expression",
		},
		{
			name:          "invalid initial term",
			input:         "+5",
			expectError:   true,
			errorContains: "invalid initial term",
		},
		{
			name:          "ends with operator",
			input:         "d20+",
			expectError:   true,
			errorContains: "expression ends with an operator",
		},
		{
			name:          "invalid operator",
			input:         "d20*5",
			expectError:   true,
			errorContains: "expected '+' or '-'",
		},
		{
			name:          "missing operator",
			input:         "d20 5",
			expectError:   true,
			errorContains: "expected '+' or '-'",
		},
		{
			name:          "malformed dice (no sides)",
			input:         "3d",
			expectError:   true,
			errorContains: "expected '+' or '-'", // counterintuitive, but correct: 3d parses 3 as constant since format (Number?)d(Number) is not satisfied by it.
		},
		{
			name:          "malformed literal",
			input:         "potato",
			expectError:   true,
			errorContains: "expected dice roll (e.g., 'd6') or modifier (e.g., '5')",
		},
		{
			name:          "double operator",
			input:         "d6++5",
			expectError:   true,
			errorContains: "expression ends with an operator",
		},
		{
			name:          "dice with zero sides",
			input:         "2d0",
			expectError:   true,
			errorContains: "dice must have more than 0 sides",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser2 := parser.NewParser(tc.input)
			got, err := parser2.Parse()

			if tc.expectError {
				if err == nil {
					t.Fatalf("Parse() expected an error, but got nil")
				}

				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Parse() error = %q, want error containing %q", err, tc.errorContains)
				}
			} else {
				if err != nil {
					t.Fatalf("Parse() returned an unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tc.expectedExpression) {
					t.Errorf("Parse() got = %+v, want %+v", got, tc.expectedExpression)
				}
			}
		})
	}
}
