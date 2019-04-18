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
	//save records the result of an action
	save(name string, err error)

	//Do executes an action as part of the sequence
	Do(name string, action Action) Sequence

	//Catch handles the consequence of a failed sequence
	Catch(consequence Consequence) Sequence

	//Then proceeds to execute a function once a sequence is marked as success
	Then(then func()) Sequence

	//Error returns the error, if any
	Error() error

	//LastAction returns the last executed action including the one that caused the error, if any
	LastAction() string
}

//unpanic is used for panic recovery without requiring it to be added to each Sequence implementation
func unpanic(name string, s Sequence) {
	if r := recover(); r != nil {
		s.save(name, ErrorPanic)
	}
}

