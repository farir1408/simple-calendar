package postgres

import (
	"context"
	"database/sql"

	"github.com/farir1408/simple-calendar/internal/pkg/types"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	// Register postgres driver.
	_ "github.com/lib/pq"
)

// Storage ...
type Storage struct {
	ctx    context.Context
	db     *sqlx.DB
	driver string
}

const defaultDriverName = "postgres"

// Opt ...
type Opt func(storage *Storage)

// WithDBConn is override db connection in storage.
func WithDBConn(db *sqlx.DB) Opt {
	return func(storage *Storage) {
		storage.db = db
	}
}

// NewFromDSN ...
func NewFromDSN(ctx context.Context, driver, dsn string, opts ...Opt) (*Storage, error) {
	if driver != "" {
		driver = defaultDriverName
	}
	s := &Storage{
		ctx:    ctx,
		driver: driver,
	}

	db, err := sqlx.Connect(s.driver, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't init postgresql connection")
	}

	s.db = db
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// SaveEvent ...
func (p *Storage) SaveEvent(ctx context.Context, event types.Event) (uint64, error) {
	result, err := p.db.QueryContext(ctx, "INSERT INTO events (title, description, start_time, duration, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", event.Title, event.Description, event.Start, event.Duration, event.UserID)
	if err != nil {
		return 0, errors.Wrap(err, "can't save event")
	}

	if !result.Next() {
		return 0, errors.New("can't save event")
	}

	var lastID uint64
	err = result.Scan(&lastID)
	if err != nil {
		return 0, errors.Wrap(err, "can't get new eventID")
	}

	return lastID, nil
}

// GetEvent ...
func (p *Storage) GetEvent(ctx context.Context, id uint64) (types.Event, error) {
	event := &DBEvent{}
	err := p.db.GetContext(ctx, event, "SELECT id, title, description, start_time, duration, user_id FROM events WHERE id=$1 AND NOT is_deleted", id)
	if err == sql.ErrNoRows {
		return types.Event{}, errors.Wrap(err, "not found")
	}
	if err != nil {
		return types.Event{}, errors.Wrap(err, "can't get event")
	}
	return convertToCalendarEvent(*event), nil
}

// UpdateEvent ...
func (p *Storage) UpdateEvent(ctx context.Context, event types.Event) error {
	res, err := p.db.ExecContext(ctx, "UPDATE events SET title=$1, description=$2, start_time=$3, duration=$4, user_id=$5 WHERE id=$6",
		event.Title,
		event.Description,
		event.Start,
		event.Duration,
		event.UserID,
		event.ID,
	)
	if err != nil {
		return errors.Wrapf(err, "can't update event: %d", event.ID)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "not updated event: %d", event.ID)
	}

	if rows == 0 {
		return errors.New("event did not update")
	}

	return nil
}

// DeleteEvent ...
func (p *Storage) DeleteEvent(ctx context.Context, id uint64) error {
	res, err := p.db.ExecContext(ctx, "UPDATE events SET is_deleted=true WHERE id=$1", id)
	if err != nil {
		return errors.Wrapf(err, "can't delete event: %d", id)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "not delete event: %d", id)
	}

	if rows == 0 {
		return errors.New("event did not delete")
	}
	return nil
}

// Close ...
func (p *Storage) Close() error {
	return p.db.Close()
}

func convertToCalendarEvent(event DBEvent) types.Event {
	return types.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Start:       event.Start,
		Duration:    event.Duration,
	}
}
