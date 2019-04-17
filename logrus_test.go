package sequitur_test

import (
	"testing"

	"github.com/fluxynet/sequitur"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var (
	logRecover = []logrus.Entry{
		{Message: "starting: letter a", Level: logrus.DebugLevel},
		{Message: "letter a", Level: logrus.InfoLevel},
		{Message: "starting: misbehave", Level: logrus.DebugLevel},
		{Message: "misbehave", Level: logrus.WarnLevel, Data: logrus.Fields{
			"error": sequitur.ErrorPanic,
		}},
		{Message: "skipping: letter b", Level: logrus.DebugLevel},
	}

	logThen = []logrus.Entry{
		{Message: "starting: letter a", Level: logrus.DebugLevel},
		{Message: "letter a", Level: logrus.InfoLevel},
		{Message: "starting: letter b", Level: logrus.DebugLevel},
		{Message: "letter b", Level: logrus.InfoLevel},
		{Message: "starting: letter c", Level: logrus.DebugLevel},
		{Message: "letter c", Level: logrus.InfoLevel},
	}

	logCatch = []logrus.Entry{
		{Message: "starting: letter a", Level: logrus.DebugLevel},
		{Message: "letter a", Level: logrus.InfoLevel},
		{Message: "starting: letter b", Level: logrus.DebugLevel},
		{Message: "letter b", Level: logrus.WarnLevel, Data: logrus.Fields{"error": errFoo}},
		{Message: "skipping: letter c", Level: logrus.DebugLevel},
	}
)

func withLogrus(f func() sequitur.Sequence, logger ...*logrus.Logger) func() sequitur.Sequence {
	return func() sequitur.Sequence {
		return sequitur.WithLogrus(f(), logger...)
	}
}

type testfn func(seq sequitur.Sequence, t *testing.T) string

func logCompare(fn testfn, h *test.Hook, expected []logrus.Entry) testfn {
	return func(seq sequitur.Sequence, t *testing.T) string {
		defer h.Reset()

		r := fn(seq, t)
		if len(h.Entries) != len(expected) {
			t.Errorf(`log entries not equal in length, expected="%d", obtained="%d"`, len(expected), len(h.Entries))
			return r
		}

		for i := range h.Entries {
			if h.Entries[i].Message != expected[i].Message {
				t.Errorf(`log entry %d does not match. expected message="%s" got message="%s"`, i, expected[i].Message, h.Entries[i].Message)
				return r
			}

			if h.Entries[i].Level != expected[i].Level {
				t.Errorf(`log entry %d does not match. expected level="%s" got level="%s"`, i, expected[i].Level, h.Entries[i].Level)
				return r
			}

			if len(h.Entries[i].Data) != len(expected[i].Data) {
				t.Errorf(`log entry %d does not match. expected data.length="%d" got data.length="%d"`, i, len(h.Entries[i].Data), len(expected[i].Data))
				return r
			}

			for k, v := range h.Entries[i].Data {
				if ve, ok := h.Entries[i].Data[k]; !ok {
					t.Errorf(`log entry %d does not match. data has unexpected field="%s"`, i, k)
				} else if ve != v {
					t.Errorf(`log entry %d does not match. expected %s="%v" got %s="%v"`, i, k, ve, k, v)
				}
			}
		}

		return r
	}
}

func TestLogrusLinearSequence(t *testing.T) {
	hook := test.NewGlobal()
	logrus.SetLevel(logrus.DebugLevel)

	testcases([]testcase{
		{"recover", withLogrus(sequitur.Linear), logCompare(test_recover, hook, logRecover), outcome{
			result: "az", lastAction: "misbehave", err: sequitur.ErrorPanic,
		}},
		{"then", withLogrus(sequitur.Linear), logCompare(test_then, hook, logThen), outcome{
			result: "abcd", lastAction: "letter c", err: nil,
		}},
		{"catch", withLogrus(sequitur.Linear), logCompare(test_catch, hook, logCatch), outcome{
			result: "abz", lastAction: "letter b", err: errFoo,
		}},
	}).run(t)
}

func TestLogrusILinearSequence(t *testing.T) {
	l, h := test.NewNullLogger()
	l.SetLevel(logrus.DebugLevel)

	testcases([]testcase{
		{"recover", withLogrus(sequitur.Linear, l), logCompare(test_recover, h, logRecover), outcome{
			result: "az", lastAction: "misbehave", err: sequitur.ErrorPanic,
		}},
		{"then", withLogrus(sequitur.Linear, l), logCompare(test_then, h, logThen), outcome{
			result: "abcd", lastAction: "letter c", err: nil,
		}},
		{"catch", withLogrus(sequitur.Linear, l), logCompare(test_catch, h, logCatch), outcome{
			result: "abz", lastAction: "letter b", err: errFoo,
		}},
	}).run(t)
}

