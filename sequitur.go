package sequitur

import "errors"

var (
	//ErrorPanic indicates a panic occurred and was recovered
	ErrorPanic = errors.New("panic")
)

//Action is part of a sequence
type Action func() error

//Consequence is the result of a sequence
type Consequence func(string, error)

//Sequence is a series of actions
type Sequence interface {
	recover()
	Do(name string, action Action) Sequence
	Catch(consequence Consequence) Sequence
	Then(then func()) Sequence
}
