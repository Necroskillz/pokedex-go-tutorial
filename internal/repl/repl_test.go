package repl

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: "pikachu", expected: []string{"pikachu"}},
		{input: "PiKaChu", expected: []string{"pikachu"}},
		{input: "", expected: []string{}},
		{input: "   pikachu   CHARMANDER     buLBasaUR   ", expected: []string{"pikachu", "charmander", "bulbasaur"}},
	}

	for _, test := range cases {
		got := CleanInput(test.input)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("cleanInput(%q) = %v; expected %v", test.input, got, test.expected)
		}
	}
}
