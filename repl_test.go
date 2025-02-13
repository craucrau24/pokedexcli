package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello   world   ",
			expected: []string{"hello", "world"},
		}, {
			input:    "All   your bases    are  belong to   us  ",
			expected: []string{"all", "your", "bases", "are", "belong", "to", "us"},
		}, {
			input:    "  Pain  is     the mindkiller",
			expected: []string{"pain", "is", "the", "mindkiller"},
		}, {
			input:    "There    is  no   spoon",
			expected: []string{"there", "is", "no", "spoon"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths mismatch %v<->%v\n", actual, c.expected)
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("item %d mismatch %v<->%v\n", i, actual[i], c.expected[i])
			}
		}
	}
}
