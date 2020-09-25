package app

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/farir1408/simple-calendar/internal/pkg/statik"

	"go.uber.org/zap"

	calendarpb "github.com/farir1408/simple-calendar/pkg/api/calendar"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"

	"github.com/rakyll/statik/fs"

	"github.com/farir1408/simple-calendar/internal/pkg/config"
	// Register static files content
	_ "github.com/farir1408/simple-calendar/internal/pkg/statik"

	"github.com/farir1408/simple-calendar/internal/pkg/types"
)

// CalendarBL calendar business logic.
type CalendarBL interface {
	io.Closer
	CreateEvent(ctx context.Context, title, description string, userID int64, start time.Time, duration time.Duration) (uint64, error)
	GetEventByID(ctx context.Context, id uint64) (types.Event, error)
	UpdateEvent(ctx context.Context, title, description string, userID int64, start time.Time, duration time.Duration) error
	DeleteEvent(ctx context.Context, id uint64) error
}

// Implementation ...
type Implementation struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *zap.Logger
	wg     sync.WaitGroup

	logic      CalendarBL
	grpcURL    string
	debugURL   string
	grpcServer *grpc.Server
	httpServer *http.Server
}

// NewApp ...
func NewApp(cfg *config.AppConfig, calendar CalendarBL, logger *zap.Logger) *Implementation {
	ctx, cancel := context.WithCancel(context.Background())
	i := &Implementation{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		logic:  calendar,
	}
	i.grpcURL = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	i.debugURL = fmt.Sprintf("%s:%d", cfg.Host, cfg.DebugPort)

	return i
}

// Start start working http and grpc server.
func (i *Implementation) Start() error {
	lis, err := net.Listen("tcp", i.grpcURL)
	if err != nil {
		return err
	}

	i.grpcServer = grpc.NewServer()
	calendarpb.RegisterCalendarServer(i.grpcServer, i)
	errs := make(chan error, 2)
	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		var grpcErr error
		grpcErr = i.grpcServer.Serve(lis)
		if grpcErr == grpc.ErrServerStopped {
			return
		}

		errs <- grpcErr
	}()
	i.logger.Info(fmt.Sprintf("start grpc service: %s", i.grpcURL))

	conn, err := grpc.DialContext(i.ctx, i.grpcURL, grpc.WithInsecure())
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux()
	err = calendarpb.RegisterCalendarHandler(i.ctx, gwmux, conn)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	mime.AddExtensionType(".svg", "image/svg+xml") //nolint:G104
	statikFS, err := fs.NewWithNamespace(statik.OpenapiDefinition)
	if err != nil {
		return err
	}
	sh := http.StripPrefix("/swaggerui/", http.FileServer(statikFS))
	mux.Handle("/swaggerui/", sh)

	i.httpServer = &http.Server{
		Addr:    i.debugURL,
		Handler: mux,
	}
	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		var httpErr error
		httpErr = i.httpServer.ListenAndServe()
		if httpErr == http.ErrServerClosed {
			return
		}
		errs <- httpErr
	}()
	i.logger.Info(fmt.Sprintf("start debug service: %s", i.debugURL))

	var resultErr error
	i.wg.Wait()
	close(errs)
	for err := range errs {
		resultErr = multierr.Append(resultErr, err)
	}

	return resultErr
}

// Close stop http and grpc service.
func (i *Implementation) Close() (err error) {
	i.cancel()
	i.logger.Info("shutdown http service...")
	err = multierr.Append(err, i.httpServer.Shutdown(i.ctx))
	i.logger.Info("shutdown grpc service...")
	i.grpcServer.Stop()
	err = multierr.Append(err, i.logic.Close())
	return err
}
