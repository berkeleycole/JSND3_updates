package main

// Race holds all of the information about the race
type Race struct {
	Cars    []*Car
	Track   *Track
	Results *RaceResults
}

// RaceOpt allows the configuration of a new Race
type RaceOpt interface {
	Apply(*Race) error
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

// RaceResults represent the results of a race
// These results can be a finised race, an in-progress race
// or a race that has yet to be started.
type RaceResults struct {
	Status    RaceStatus `json:"status"`
	Positions Position   `json:"positions"`
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

// Position wraps a car and determines the position
// of a car in a race.
type Position struct {
	Car
	Position int `json:"position"`
}
