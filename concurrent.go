package sequitur

//ConcurrentSequence is a group of actions that happen concurrently
type ConcurrentSequence struct {
	//Error represents any error that has stopped the execution of a sequence
	Error error

	//LastAction represents the last executed action on the sequence
	LastAction string

	//pending indicates that the sequence is still pending, will still accept any do, catch or then
	pending bool
}

func (s *ConcurrentSequence) recover() {

}

//Do queues an action for execution
func (s *ConcurrentSequence) Do(name string, action Action) Sequence {
	return s
}

//Catch sets up the consequence to catch a failing sequence
func (s *ConcurrentSequence) Catch(consequence Consequence) Sequence {
	return s
}

//Then sets up the result of the sequence, and also triggers the start of the sequence
func (s *ConcurrentSequence) Then(then func()) Sequence {
	return s
}

