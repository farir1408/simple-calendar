package types

import "time"

// Event ...
type Event struct {
	ID          uint64
	Title       string
	Description string
	Start       time.Time
	Duration    time.Duration
	UserID      int64
}
