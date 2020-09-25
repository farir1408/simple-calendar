//+build integration

package postgres

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/farir1408/simple-calendar/internal/pkg/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/suite"
)

const (
	testDBDSN      = "postgres://calendar-api-user:calendar-api-password@localhost:5432/calendar-api-db?sslmode=disable"
	testDriverName = "pgx"
)

type OLTPTestSuite struct {
	suite.Suite
	ctx  context.Context
	conn *sqlx.DB
	db   *Storage
}

func (suite *OLTPTestSuite) SetupTest() {
	suite.ctx = context.Background()
	conn, err := sqlx.Connect(testDriverName, testDBDSN)
	suite.Require().NoError(err)
	db, err := NewFromDSN(suite.ctx, testDriverName, testDBDSN, WithDBConn(conn))
	suite.Require().NoError(err)
	suite.db = db
	suite.conn = conn
	err = suite.conn.Ping()
	suite.Require().NoError(err)
}

func (suite *OLTPTestSuite) TearDownTest() {
	suite.cleanTables(tableNameEvents)
}

func (suite *OLTPTestSuite) cleanTables(tableNames ...string) {
	_, _ = suite.conn.ExecContext(suite.ctx, fmt.Sprintf("TRUNCATE TABLES %s CASACDE", strings.Join(tableNames, ", ")))
}

func (suite *OLTPTestSuite) Test_CreateEvent() {
	startTime := time.Now()
	event := types.Event{
		Title:       "test",
		Description: "description",
		Start:       startTime,
		Duration:    time.Hour,
		UserID:      100,
	}

	id, err := suite.db.SaveEvent(suite.ctx, event)
	suite.Require().NoError(err)
	suite.Require().NotEqual(0, id)
}

func (suite *OLTPTestSuite) Test_GetEvent() {
	startTime := time.Now()
	event := types.Event{
		Title:       "test",
		Description: "description",
		Start:       startTime,
		Duration:    time.Hour,
		UserID:      100,
	}
	id, err := suite.db.SaveEvent(suite.ctx, event)
	suite.Require().NoError(err)

	location, err := time.LoadLocation("Europe/Moscow")
	suite.Require().NoError(err)

	dbEvent, err := suite.db.GetEvent(suite.ctx, id)
	suite.Require().NoError(err)
	suite.Require().Equal(startTime.Format(time.RFC1123), dbEvent.Start.In(location).Format(time.RFC1123))
	suite.Require().Equal(time.Hour, dbEvent.Duration)
}

func TestOLTPSuite(t *testing.T) {
	suite.Run(t, new(OLTPTestSuite))
}
