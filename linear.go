package sequitur

//LinearSequence is a group of actions
type LinearSequence struct {
	//Error represents any error that has stopped the execution of a sequence
	Error error

	//LastAction represents the last executed action on the sequence
	LastAction string

	//pending indicates that the sequence is still pending, will still accept any do, catch or then
	pending bool
}

//Linear returns a new linear sequence
func Linear() *LinearSequence {
	return &LinearSequence{pending: true}
}

func (s *LinearSequence) save(name string, err error) {
	s.LastAction = name
	s.Error = err
}

//Do executes an action as part of a sequence
func (s *LinearSequence) Do(name string, action Action) Sequence {
	if s.pending && s.Error == nil {
		var err error
		defer unpanic(name, s)
		defer s.save(name, err)
		err = action()
	}

	return s
}

//Catch executes a consequence if an error has occurred as part of a sequence
func (s *LinearSequence) Catch(consequence Consequence) Sequence {
	if s.pending && s.Error != nil {
		consequence(s.LastAction, s.Error)
		s.pending = false
	}

	return s
}

//Then executes a function if no error has occurred
func (s *LinearSequence) Then(then func()) Sequence {
	if s.pending && s.Error == nil {
		then()
		s.pending = false
	}

	return s
}
