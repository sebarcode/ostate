package ostate_test

import (
	"testing"

	"github.com/sebarcode/ostate"
	"github.com/smartystreets/goconvey/convey"
)

func TestOstate(t *testing.T) {
	convey.Convey("test ostate", t, func() {

		req := &sampleData{
			ID:     "Req1",
			Name:   "Request 1",
			Value:  1000,
			Status: "Draft",
		}

		sm := ostate.NewStateEngine(&ostate.StateEngineOpts{
			RegulateFlow: true,
		})

		err := sm.Load(req)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("submit", func() {
			err = sm.ChangeState(req, "Submitted", nil)
			convey.So(err, convey.ShouldBeNil)
			convey.So(req.State(), convey.ShouldEqual, "Submitted")
		})
	})
}

type sampleData struct {
	ID     string
	Name   string
	Value  int
	Status string
}

func (s *sampleData) Stats() (id string, title string, state string, err error) {
	return s.ID, s.Name, s.Status, nil
}

func (s *sampleData) ChangeState(newState string, payload interface{}) error {
	s.Status = newState
	return nil
}

func (s *sampleData) State() string {
	return s.Status
}
