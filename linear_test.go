package sequitur_test

import (
	"errors"
	"testing"

	"github.com/fluxynet/sequitur"
)

func TestLinearSequence_recover(t *testing.T) {
	var (
		seq    = sequitur.Linear()
		result string
		gotErr error
	)

	seq.Do("letter a", func() error {
		result += "a"
		return nil
	})

	seq.Do("misbehave", func() error {
		panic("foobar")
	})

	seq.Catch(func(name string, err error) {
		gotErr = err
		result = name
	})

	if result != "misbehave" {
		t.Errorf(`expected name="misbehave", obtained="%s"`, result)
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

	defer func() {
		if result != "abcd" {
			t.Errorf(`expected "abcd", obtained="%s"\n`, result)
		}
	}()

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
}

func TestLinearSequence_DoCatch(t *testing.T) {
	var (
		seq     = sequitur.Linear()
		result  string
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
		result = name
		errGot = err
	})

	seq.Then(func() {
		result += "d"
	})

	if result != "letter b" {
		t.Errorf(`expected name="letter b", obtained="%s"\n`, result)
	}

	if errGot != errWant {
		t.Errorf(`expected error=[%v], obtained=[%v]`, errWant, errGot)
	}
}
