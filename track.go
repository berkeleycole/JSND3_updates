package main

// MaxTurnDegree represents the absolute value of the maximum
// degree that a single segment of a track can be.
const MaxTurnDegree = 90

// A Track is represented as an array of segments
type Track struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Segments []float64 `json:"segments"`
}

// A Segment represents a single tile of the race track
// It has a value for the degree of the turn of that tile.
// The max for this value is 90 when expressed as a degree
type Segment struct {
	Degree  int     `json:"degree"`
	Percent float64 `json:"percent"`
}
