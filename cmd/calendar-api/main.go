package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"syscall"

	"github.com/farir1408/simple-calendar/internal/pkg/storage/postgres"

	"github.com/farir1408/simple-calendar/pkg/closer"

	"github.com/farir1408/simple-calendar/internal/pkg/calendar"

	"github.com/farir1408/simple-calendar/internal/app"
	"github.com/farir1408/simple-calendar/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	cfg := config.New()
	lvl := zap.NewAtomicLevel()
	lvl.SetLevel(setLogLvl(cfg))
	logger, err := zap.NewProduction(zap.IncreaseLevel(lvl))
	if err != nil {
		log.Fatal("can't init logger", err)
	}
	//ctx := ctxzap.ToContext(ctx, logger)
	repository, err := postgres.NewFromDSN(ctx, "postgres", createDSN(cfg.DatabaseAddr, cfg.DatabaseDBName, cfg.DatabaseUser, cfg.DatabasePassword))
	if err != nil {
		log.Fatal(err)
	}
	logic := calendar.NewCalendar(ctx, repository)

	service := app.NewApp(cfg, logic, logger)

	c := closer.New(syscall.SIGINT, syscall.SIGTERM)
	c.Bind(service.Close)

	fmt.Println("Start...")
	if err := service.Start(); err != nil {
		log.Fatal(err)
	}
}

//TODO: ...
func setLogLvl(cfg *config.AppConfig) zapcore.Level {
	if cfg == nil {
		fmt.Println("LogLevel is not set.")
		return zapcore.InfoLevel
	}

	switch strings.ToLower(cfg.LogLvl) {
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func createDSN(addr, dbName, user, password string) string {
	dbURL := &url.URL{
		Scheme:   "postgres",
		Host:     addr,
		User:     url.UserPassword(user, password),
		Path:     dbName,
		RawQuery: "sslmode=disable",
	}
	fmt.Println("database: ", dbURL.String())
	return dbURL.String()
}
