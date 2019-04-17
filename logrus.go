package sequitur

import (
	"github.com/sirupsen/logrus"
)

type sequenceWithLogrus struct {
	Sequence
}

type sequenceWithLogrusI struct {
	Sequence
	logger *logrus.Logger
}

//WithLogrus returns a wrapped sequence with logrus logging features
//will use global logrus logger unless a non-nil logger instance is passed as second argument
func WithLogrus(s Sequence, l ...*logrus.Logger) Sequence {
	if len(l) == 0 || l[0] == nil {
		return &sequenceWithLogrus{s}
	}

	return &sequenceWithLogrusI{s, l[0]}
}
func (s sequenceWithLogrus) Do(name string, action Action) Sequence {
	if err := s.Sequence.Error(); err != nil {
		logrus.Debug("skipping: " + name)
		return s
	}

	logrus.Debug("starting: " + name)
	s.Sequence.Do(name, action)

	if err := s.Sequence.Error(); err == nil {
		logrus.Info(name)
	} else {
		logrus.WithError(err).Warn(name)
	}

	return s
}

func (s sequenceWithLogrusI) Do(name string, action Action) Sequence {
	if err := s.Sequence.Error(); err != nil {
		s.logger.Debug("skipping: " + name)
		return s
	}

	s.logger.Debug("starting: " + name)
	s.Sequence.Do(name, action)

	if err := s.Sequence.Error(); err == nil {
		s.logger.Info(name)
	} else {
		s.logger.WithError(err).Warn(name)
	}

	return s
}
