package sequitur

import (
	"context"
)

type sequenceWithContext struct {
	Sequence
	Ctx context.Context
}

//WithContext returns a wrapped sequence with context checking
func WithContext(ctx context.Context, s Sequence) Sequence {
	return &sequenceWithContext{
		Sequence: s,
		Ctx:      ctx,
	}
}

//Do wraps the normal Do with a context check
func (s *sequenceWithContext) Do(name string, action Action) Sequence {
	select {
	case <-s.Ctx.Done():
		s.save(name, s.Ctx.Err())
	default:
		s.Sequence.Do(name, action)
	}

	return s
}

func (s *sequenceWithContext) Then(then func()) Sequence {
	s.Sequence.Then(then)
	return s
}

func (s *sequenceWithContext) Catch(consequence Consequence) Sequence {
	s.Sequence.Catch(consequence)
	return s
}
