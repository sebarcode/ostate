package ostate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebarcode/codekit"
	"github.com/sebarcode/ostate"
	"github.com/smartystreets/goconvey/convey"
)

func TestOstate(t *testing.T) {
	convey.Convey("test ostate", t, func() {

		req := &sampleData{
			ID:          "Req1",
			Name:        "Request 1",
			Value:       1000,
			CompletedBy: "",
			Status:      "Draft",
		}

		sm := ostate.NewStateEngine(&ostate.StateEngineOpts{
			RegulateFlow: true,
			Flow: map[string][]string{
				"Draft":     {"Submitted"},
				"Submitted": {"Working", "Rejected"},
				"Working":   {"Done", "Error"},
				"Rejected":  {"Error"},
			},
		})

		err := sm.Load(req)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("submit", func() {
			err = sm.ChangeState(req, "Submitted", nil)
			convey.So(err, convey.ShouldBeNil)
			convey.So(req.State(), convey.ShouldEqual, "Submitted")

			convey.Convey("invalid state", func() {
				err = sm.ChangeState(req, "Done", nil)
				convey.So(err, convey.ShouldNotBeNil)
				convey.Print("  ", req.State())

				convey.Convey("change to working", func() {
					err = sm.ChangeState(req, "Working", nil)
					convey.So(err, convey.ShouldBeNil)

					convey.Convey("change to done", func() {
						err = sm.ChangeState(req, "Done", codekit.M{}.Set("Operator", "Samsul"))
						convey.So(err, convey.ShouldBeNil)
						convey.So(int(req.Duration), convey.ShouldBeGreaterThan, 0)
						convey.Printf(codekit.JsonString(codekit.M{}.Set("operated", req.CompletedBy).Set("duration", req.Duration.String())))
					})
				})
			})
		})
	})
}

type sampleData struct {
	ID                 string
	Name               string
	Value              int
	Started, Completed *time.Time
	Duration           time.Duration
	CompletedBy        string
	Status             string
}

func (s *sampleData) Stats() (id string, title string, state string, err error) {
	return s.ID, s.Name, s.Status, nil
}

func (s *sampleData) ChangeState(newState string, payload interface{}) error {
	s.Status = newState

	switch newState {
	case "Working":
		completed := time.Now()
		s.Started = &completed

	case "Done":
		completed := time.Now()
		s.Completed = &completed
		s.Duration = s.Completed.Sub(*s.Started)
		m, ok := payload.(codekit.M)
		if !ok {
			return fmt.Errorf("invalid payload, %t", payload)
		}
		s.CompletedBy = m.GetString("Operator")
	}

	return nil
}

func (s *sampleData) State() string {
	return s.Status
}
