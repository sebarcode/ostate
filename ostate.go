package ostate

import (
	"fmt"
	"slices"
)

type StateObject interface {
	Stats() (id, title, state string, err error)
	ChangeState(newState string, payload interface{}) error
	State() string
}

type StateEngineOpts struct {
	RegulateFlow bool
	Flow         map[string][]string
}

type stateEngine struct {
	opts StateEngineOpts
}

func NewStateEngine(opts *StateEngineOpts) *stateEngine {
	s := new(stateEngine)
	if opts != nil {
		s.opts = *opts
	}
	return s
}

func (s *stateEngine) Load(obj StateObject) error {
	return nil
}

func (s *stateEngine) ChangeState(obj StateObject, newState string, payload interface{}) error {
	if s.opts.RegulateFlow {
		// validate flow
		allowedStates, ok := s.opts.Flow[obj.State()]
		if ok {
			if slices.Index(allowedStates, newState) < 0 {
				return fmt.Errorf("invalid new state: %s", newState)
			}
		}
	}
	return obj.ChangeState(newState, payload)
}
