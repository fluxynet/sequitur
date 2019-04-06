package sequitur

//Sequence is a group of actions
type Sequence struct {
	//Error represents any error that has stopped the execution of a sequence
	Error error

	//LastAction represents the last executed action on the sequence
	LastAction string
}

//Action is part of a sequence
type Action func() error

//Consequence is the result of a sequence
type Consequence func(string, error)

//Do executes an action as part of a sequence
func (s Sequence) Do(name string, action Action) Sequence {
	if s.Error == nil {
		s.LastAction = name
		s.Error = action()
	}

	return s
}

//Catch executes a consequence if an error has occurred as part of a sequence
func (s Sequence) Catch(consequence Consequence) Sequence {
	if s.Error != nil {
		consequence(s.LastAction, s.Error)
	}

	return s
}

//Then executes a function if no error has occurred
func (s Sequence) Then(then func()) Sequence {
	if s.Error == nil {
		then()
	}

	return s
}
