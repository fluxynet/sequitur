package sequitur

//LinearSequence is a group of actions
type LinearSequence struct {
	//error represents any error that has stopped the execution of a sequence
	err error

	//lastAction represents the last executed action on the sequence
	lastAction string

	//pending indicates that the sequence is still pending, will still accept any do, catch or then
	pending bool
}

//Linear returns a new linear sequence
func Linear() Sequence {
	return &LinearSequence{pending: true}
}

func (s *LinearSequence) save(name string, err error) {
	s.lastAction = name
	s.err = err
}

//Error returns the error, if any
func (s *LinearSequence) Error() error {
	return s.err
}

//LastAction returns the last executed action including the one that caused the error, if any
func (s *LinearSequence) LastAction() string {
	return s.lastAction
}

//Do executes an action as part of a sequence
func (s *LinearSequence) Do(name string, action Action) {
	if s.pending && s.err == nil {
		var err error
		defer unpanic(name, s)
		defer func() {
			s.save(name, err)
		}()
		err = action()
	}
}

//Catch executes a consequence if an error has occurred as part of a sequence
func (s *LinearSequence) Catch(consequence Consequence) {
	if s.pending && s.err != nil {
		consequence(s.lastAction, s.err)
		s.pending = false
	}
}

//Then executes a function if no error has occurred
func (s *LinearSequence) Then(then func()) {
	if s.pending && s.err == nil {
		then()
		s.pending = false
	}
}
