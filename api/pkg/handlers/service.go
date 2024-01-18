package handlers

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/routes"
	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

type EchoContext interface {
	echo.Context
}

type EchoServer interface {
	Start(string) error
	Shutdown(ctx context.Context) error
}

type Service struct {
	ctx    context.Context
	opts   *Options
	logger logger.Logger
	server EchoServer

	Services Services
}

// Dependencies - dependencies for Service constructor
type Dependencies struct {
	Logger        logger.Logger
	EchoServer    EchoServer
	MessageBroker MessageBroker
	GormDB        *gorm.DB
	Services      Services
}

type RabbitMqService struct {
	HostPort string
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

type Options struct {
	Path                string
	Port                int
	ShutdownGracePeriod time.Duration
}

type Services struct {
	OpenGraphSvc OpenGraphService
}

// GetFlagSet returns flag set for Options
func (o *Options) GetFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("apiOptions", pflag.ExitOnError)
	flags.StringVar(&o.Path, "path", o.Path, "path to serve API on")
	flags.IntVar(&o.Port, "port", o.Port, "port to serve API on")
	flags.DurationVar(&o.ShutdownGracePeriod, "shutdown-grace-period", o.ShutdownGracePeriod, "shutdown grace period")
	return flags
}

// NewService - constructor for Service
func NewService(ctx context.Context, opts *Options, deps *Dependencies) (*Service, error) {
	svc := &Service{
		ctx:      ctx,
		opts:     opts,
		logger:   deps.Logger,
		server:   deps.EchoServer,
		Services: deps.Services,
	}
	svc.server = svc.createServer()
	return svc, nil
}

// Start starts the API
func (svc *Service) Start() {
	go func() {
		addr := fmt.Sprintf(":%d", svc.opts.Port)
		if err := svc.server.Start(addr); err != nil {
			logger.Println(err)
		}
	}()
}

// Close closes the API
func (svc *Service) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.opts.ShutdownGracePeriod)
	defer cancel()
	return svc.server.Shutdown(ctx)
}

func (svc *Service) createServer() EchoServer {
	server := echo.New()
	server.Use(middleware.CORS())
	server.JSONSerializer = &jsonSerializer{}
	// server.Use(svc.AuthzMiddleware)
	apiGroup := server.Group("")
	routes.RegisterHandlersWithBaseURL(apiGroup, svc, svc.opts.Path)
	return server
}
