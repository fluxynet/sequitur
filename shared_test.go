package sequitur_test

import (
	"errors"
	"testing"
	"time"

	"github.com/fluxynet/sequitur"
)

var (
	errFoo = errors.New("foo")
)

type outcome struct {
	result     string
	lastAction string
	err        error
}

type testcase struct {
	name     string
	seq      func() sequitur.Sequence
	test     func(seq sequitur.Sequence, t *testing.T) string
	expected outcome
}

func (tt testcase) run(t *testing.T) {
	seq := tt.seq()
	result := tt.test(seq, t)

	if tt.expected.result != result {
		t.Errorf(`expected result="%s", obtained="%s"`, tt.expected.result, result)
	}

	if seq.LastAction() != tt.expected.lastAction {
		t.Errorf(`expected lastAction="%s", obtained="%s"`, tt.expected.lastAction, seq.LastAction())
	}

	if tt.expected.err != seq.Error() {
		t.Errorf(`expected err="%s", obtained="%s"`, tt.expected.err, seq.Error())
	}
}

type testcases []testcase

func (tc testcases) run(t *testing.T) {
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}

func test_recover(seq sequitur.Sequence, t *testing.T) string {
	var result string

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("misbehave", func() error {
		panic("foobar")
	})

	seq.Do("letter b", func() error {
		result += "b"
		return nil
	})

	seq.Catch(func(name string, err error) {
		result += "z"
	})

	return result
}

func test_then(seq sequitur.Sequence, t *testing.T) string {
	var result string

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("letter b", func() error {
		result += "b"
		return nil
	})

	seq.Do("letter c", func() error {
		result += "c"
		return nil
	})

	seq.Catch(func(name string, err error) {
		result += "z"
	})

	seq.Then(func() {
		result += "d"
	})

	seq.Then(func() {
		result += "e"
	})

	return result
}

func test_then_delayed(seq sequitur.Sequence, t *testing.T) string {
	var result string

	seq.Do("letter a", func() error {
		time.Sleep(time.Second * 10)
		result += "a"
		return nil
	})

	seq.Do("letter b", func() error {
		result += "b"
		return nil
	})

	seq.Do("letter c", func() error {
		result += "c"
		return nil
	})

	seq.Catch(func(name string, err error) {
		result += "z"
	})

	seq.Then(func() {
		result += "d"
	})

	seq.Then(func() {
		result += "e"
	})

	return result
}

func test_catch(seq sequitur.Sequence, t *testing.T) string {
	var result string

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("letter b", func() error {
		result += "b"
		return errFoo
	})

	seq.Do("letter c", func() error {
		result += "c"
		return nil
	})

	seq.Catch(func(name string, err error) {
		result += "z"
	})

	seq.Then(func() {
		result += "d"
	})

	return result
}
