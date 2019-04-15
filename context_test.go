package sequitur_test

import (
	"context"
	"testing"

	"github.com/fluxynet/sequitur"
)

func withCtx(ctx context.Context, f func() sequitur.Sequence) func() sequitur.Sequence {
	return func() sequitur.Sequence {
		return sequitur.WithContext(ctx, f())
	}
}

func TestContextLinearSequence(t *testing.T) {
	ctxBG := context.Background()
	ctxC, cancel := context.WithCancel(context.Background())

	cancel()

	testcases([]testcase{
		{"recover", withCtx(ctxBG, sequitur.Linear), test_recover, outcome{
			result: "az", lastAction: "misbehave", err: sequitur.ErrorPanic,
		}},
		{"then", withCtx(ctxBG, sequitur.Linear), test_then, outcome{
			result: "abcd", lastAction: "letter c", err: nil,
		}},
		{"catch", withCtx(ctxBG, sequitur.Linear), test_catch, outcome{
			result: "abz", lastAction: "letter b", err: errFoo,
		}},
		{"then_delayed", withCtx(ctxC, sequitur.Linear), test_then_delayed, outcome{
			result: "z", lastAction: "letter a", err: context.Canceled,
		}},
	}).run(t)
}
