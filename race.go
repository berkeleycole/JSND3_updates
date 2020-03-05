package main

import (
	"time"
)

// Race holds all of the information about the race
type Race struct {
	Track                *Track
	PlayerID             int
	clicksSinceLastCheck int
	Cars                 []*Car
	Results              *RaceResults
	actionCh             chan int
	tickerCh             chan time.Time
	closeCh              chan struct{}
}

// NewRace constructs a Race from the provided options
func NewRace(opts ...RaceOpt) (*Race, error) {
	r := Race{}

	for _, opt := range opts {
		if err := opt.Apply(&r); err != nil {
			return nil, err
		}
	}

	return &r, nil
}

/*

Actions:
	- Start
	- Finish
	- Refresh (passing of time) - computes the location of each car on the track based on speed
	- HandleAction - handles the "click" from the user

*/

// Start starts a race
func (r *Race) Start() error {
	r.Results.Status = InProgress

	go r.loop()

	return nil
}

// Finish marks the end of a race
func (r *Race) Finish() error {
	close(r.closeCh)
	close(r.tickerCh)
	close(r.actionCh)

	return nil
}

// Refresh computes the location of cars on the track
// based on speed and acceleration
func (r *Race) Refresh() (*RaceResults, error) {
	raceFinished := true
	finishPosition := 0

	for _, car := range r.Results.Positions {
		if car.ID == r.PlayerID {
			if r.clicksSinceLastCheck > 0 {
				car.Speed = car.Speed + (car.Acceleration / 3) + r.clicksSinceLastCheck
				r.clicksSinceLastCheck = 0
			} else {
				car.Speed = car.Speed - 20
			}
		} else {
			car.Speed = car.Speed + car.Acceleration
		}

		if car.Speed > car.TopSpeed {
			car.Speed = car.TopSpeed
		}

		car.Segment = car.Segment + (car.Speed / 30)
		if car.Segment >= len(r.Track.Segments) {
			car.Segment = len(r.Track.Segments)
			finishPosition++
			car.FinalPosition = finishPosition
		} else {
			raceFinished = false
		}
	}

	if raceFinished {
		r.Finish()
	}

	return r.Results, nil
}

func (r *Race) loop() {
	for {
		select {
		case <-r.tickerCh:
			r.Refresh()
		case <-r.actionCh:
			r.clicksSinceLastCheck++
		case <-r.closeCh:
			r.Results.Status = Finished
			return
		}

		time.Sleep(time.Second * 1)
	}
}

// Accelerate handles the action from the user to increase the speed
// Acceleration sets the cap for this.
func (r *Race) Accelerate() error {
	r.actionCh <- 0

	return nil
}

// RaceOpt allows the configuration of a new Race
type RaceOpt interface {
	Apply(*Race) error
}

// RaceResults represent the results of a race
// These results can be a finised race, an in-progress race
// or a race that has yet to be started.
type RaceResults struct {
	Status    RaceStatus     `json:"status"`
	Positions []*CarPosition `json:"positions"`
}

// RaceStatus is the status of the race.
// It is eaither unstarted, in-progress, or finished
type RaceStatus string

// Values for RaceStatus
const (
	Unstarted  = "unstarted"
	InProgress = "in-progress"
	Finished   = "finished"
)

// CarPosition wraps a car and determines the position
// of a car in a race.
type CarPosition struct {
	Car
	FinalPosition int `json:"final_position,omitempty"`
	Speed         int `json:"speed"`
	Segment       int `json:"segment"`
}
