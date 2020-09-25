package postgres

import "time"

// DBEvent ...
type DBEvent struct {
	ID          uint64        `db:"id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Start       time.Time     `db:"start_time"`
	Duration    time.Duration `db:"duration"`
	UserID      uint64        `db:"user_id"`
	IsDeleted   bool          `db:"is_deleted"`
}
