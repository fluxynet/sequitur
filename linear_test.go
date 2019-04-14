package sequitur_test

import (
	"testing"

	"github.com/fluxynet/sequitur"
)

func TestLinearSequence(t *testing.T) {
	testcases([]testcase{
		{"recover", sequitur.Linear, test_recover, outcome{
			result: "az", lastAction: "misbehave", err: sequitur.ErrorPanic,
		}},
		{"then", sequitur.Linear, test_then, outcome{
			result: "abcd", lastAction: "letter c", err: nil,
		}},
		{"catch", sequitur.Linear, test_catch, outcome{
			result: "abz", lastAction: "letter b", err: errFoo,
		}},
	}).run(t)
}
