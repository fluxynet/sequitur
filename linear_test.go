package sequitur_test

import (
	"errors"
	"testing"

	"github.com/fluxynet/sequitur"
)

func TestLinearSequence_recover(t *testing.T) {
	var (
		seq     = sequitur.Linear()
		result  string
		gotErr  error
		gotName string
	)

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
		gotErr = err
		gotName = name
	})

	if result != "a" {
		t.Errorf(`expected result="a", obtained="%s"`, result)
	}

	if gotName != "misbehave" {
		t.Errorf(`expected name="misbehave", obtained="%s"`, gotName)
	}

	if gotErr != sequitur.ErrorPanic {
		t.Errorf(`expected error=[%v], obtained=[%v]`, sequitur.ErrorPanic, gotErr)
	}
}

func TestLinearSequence_DoThen(t *testing.T) {
	var (
		seq    = sequitur.Linear()
		result string
	)

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
		result += "e"
	})

	seq.Then(func() {
		result += "d"
	})

	seq.Then(func() {
		result += "z"
	})

	if result != "abcd" {
		t.Errorf(`expected "abcd", obtained="%s"\n`, result)
	}
}

func TestLinearSequence_DoCatch(t *testing.T) {
	var (
		seq     = sequitur.Linear()
		result  string
		nameGot string
		errGot  error
		errWant = errors.New("error")
	)

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("letter b", func() error {
		result += "b"
		return errWant
	})

	seq.Do("letter c", func() error {
		result += "c"
		return nil
	})

	seq.Catch(func(name string, err error) {
		nameGot = name
		errGot = err
	})

	seq.Then(func() {
		result += "d"
	})

	if result != "ab" {
		t.Errorf(`expected result="ab", obtained="%s"\n`, result)
	}

	if nameGot != "letter b" {
		t.Errorf(`expected name="letter b", obtained="%s"\n`, nameGot)
	}

	if errGot != errWant {
		t.Errorf(`expected error=[%v], obtained=[%v]`, errWant, errGot)
	}
}
