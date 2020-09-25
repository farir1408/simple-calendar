package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/farir1408/simple-calendar/internal/pkg/types"
	"github.com/stretchr/testify/suite"
)

type InMemoryStorageTestSuite struct {
	suite.Suite
	ctx     context.Context
	storage *Storage
}

func (suite *InMemoryStorageTestSuite) SetupTest() {
	// before each test
	suite.ctx = context.Background()
	suite.storage = NewInMemoryStorage(10)
}

func (suite *InMemoryStorageTestSuite) TearDownTest() {
	// after each test
}

func (suite *InMemoryStorageTestSuite) Test_SaveEventOk() {
	event := types.Event{
		UserID:      1000,
		Title:       "test",
		Description: "test event",
		Start:       time.Now(),
		Duration:    time.Hour,
	}
	eventID, err := suite.storage.SaveEvent(suite.ctx, event)
	suite.Require().NoError(err)
	suite.Require().NotEqual(0, eventID)
}

func (suite *InMemoryStorageTestSuite) Test_GetEventOk() {
	eventStart := time.Now()
	event := types.Event{
		UserID:      1000,
		Title:       "test",
		Description: "test event",
		Start:       eventStart,
		Duration:    time.Hour,
	}
	eventID, err := suite.storage.SaveEvent(suite.ctx, event)
	suite.Require().NoError(err)
	suite.Require().NotEqual(0, eventID)

	dbEvent, err := suite.storage.GetEvent(suite.ctx, eventID)
	suite.Require().NoError(err)
	suite.Require().Equal("test", dbEvent.Title)
	suite.Require().Equal("test event", dbEvent.Description)
	suite.Require().Equal(eventStart, dbEvent.Start)
	suite.Require().Equal(time.Hour, dbEvent.Duration)
}

func (suite *InMemoryStorageTestSuite) Test_GetEventNotExist() {
	_, err := suite.storage.GetEvent(suite.ctx, 1000)
	suite.Require().Error(err)
}

func (suite *InMemoryStorageTestSuite) Test_UpdateEventOk() {
	eventStart := time.Now()
	event := types.Event{
		UserID:      1000,
		Title:       "test",
		Description: "test event",
		Start:       eventStart,
		Duration:    time.Hour,
	}
	eventID, err := suite.storage.SaveEvent(suite.ctx, event)
	suite.Require().NoError(err)
	suite.Require().NotEqual(0, eventID)

	dbEvent, err := suite.storage.GetEvent(suite.ctx, eventID)
	suite.Require().NoError(err)
	suite.Require().Equal("test", dbEvent.Title)

	dbEvent.Title = "updated title"
	err = suite.storage.UpdateEvent(suite.ctx, dbEvent)
	suite.Require().NoError(err)
	suite.Require().Equal("updated title", dbEvent.Title)
}

func (suite *InMemoryStorageTestSuite) Test_UpdateEventNotExist() {
	err := suite.storage.UpdateEvent(suite.ctx, types.Event{})
	suite.Require().Error(err)
}

func TestInMemoryStorage(t *testing.T) {
	suite.Run(t, new(InMemoryStorageTestSuite))
}
