package calendar

import (
	"context"
	"io"
	"time"

	"github.com/farir1408/simple-calendar/internal/pkg/types"
	"github.com/pkg/errors"
)

// Storage ...
type Storage interface {
	io.Closer
	SaveEvent(ctx context.Context, event types.Event) (uint64, error)
	GetEvent(ctx context.Context, id uint64) (types.Event, error)
	UpdateEvent(ctx context.Context, event types.Event) error
	DeleteEvent(ctx context.Context, id uint64) error
}

// Opt ...
type Opt func(c *Calendar)

// Calendar ...
type Calendar struct {
	ctx     context.Context
	cancel  context.CancelFunc
	storage Storage
}

// NewCalendar ...
func NewCalendar(ctx context.Context, storage Storage, opts ...Opt) *Calendar {
	ctx, cancel := context.WithCancel(ctx)
	c := &Calendar{
		ctx:     ctx,
		cancel:  cancel,
		storage: storage,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// CreateEvent ...
func (c *Calendar) CreateEvent(ctx context.Context, title, description string, userID int64, start time.Time, duration time.Duration) (uint64, error) {
	err := c.validateStartEventTime(start)
	if err != nil {
		return 0, err
	}

	err = c.validateEventDuration(duration)
	if err != nil {
		return 0, err
	}

	return c.storage.SaveEvent(ctx, types.Event{
		Title:       title,
		Description: description,
		Start:       start,
		Duration:    duration,
		UserID:      userID,
	})
}

// GetEventByID ...
func (c *Calendar) GetEventByID(ctx context.Context, id uint64) (types.Event, error) {
	return c.storage.GetEvent(ctx, id)
}

// UpdateEvent ...
func (c *Calendar) UpdateEvent(ctx context.Context, title, description string, userID int64, start time.Time, duration time.Duration) error {
	err := c.validateStartEventTime(start)
	if err != nil {
		return err
	}

	err = c.validateEventDuration(duration)
	if err != nil {
		return err
	}

	return c.storage.UpdateEvent(ctx, types.Event{
		Title:       title,
		Description: description,
		Start:       start,
		Duration:    duration,
		UserID:      userID,
	})
}

// DeleteEvent ...
func (c *Calendar) DeleteEvent(ctx context.Context, id uint64) error {
	return c.storage.DeleteEvent(ctx, id)
}

// Close ...
func (c *Calendar) Close() error {
	c.cancel()
	return c.storage.Close()
}

func (c *Calendar) validateStartEventTime(startTime time.Time) error {
	if startTime.Before(time.Now()) {
		return errors.New("start time of event before now")
	}
	return nil
}

func (c *Calendar) validateEventDuration(duration time.Duration) error {
	if duration < 0 {
		return errors.New("duration is low then zero")
	}

	if duration == 0 {
		return errors.New("empty duration")
	}

	return nil
}
