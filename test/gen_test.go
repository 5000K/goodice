package test

import (
	"testing"

	"github.com/5000K/goodice"
)

func TestNew_ErrorCases(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
	}{
		{"invalid expression", "d20*5"},
		{"ends with operator", "3d6-"},
		{"starts with operator", "+4"},
		{"empty expression", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := goodice.New(tc.expression)
			if err == nil {
				t.Errorf("New(%q) expected an error but got nil", tc.expression)
			}
		})
	}
}

// smoke test
func TestGoodice_Generate(t *testing.T) {
	expression := "2d6+5"
	roller, err := goodice.New(expression)
	if err != nil {
		t.Fatalf("New(%q) failed with unexpected error: %v", expression, err)
	}

	result, err := roller.Generate()
	if err != nil {
		t.Fatalf("Generate() returned an unexpected error: %v", err)
	}

	expectedNumParts := 2
	if len(result.Parts) != expectedNumParts {
		t.Errorf("Expected %d result parts, but got %d", expectedNumParts, len(result.Parts))
	}

	part1 := result.Parts[0]
	if part1.Type != goodice.DiceRoll {
		t.Errorf("Part 1: expected type DiceRoll, got %v", part1.Type)
	}

	expectedDiceCount := 2
	if part1.Sides != expectedDiceCount {
		t.Errorf("Part 1: expected Sides (dice count) to be %d, got %d", expectedDiceCount, part1.Sides)
	}

	if len(part1.ResultParts) != expectedDiceCount {
		t.Errorf("Part 1: expected %d individual die roll results, got %d", expectedDiceCount, len(part1.ResultParts))
	}
}
