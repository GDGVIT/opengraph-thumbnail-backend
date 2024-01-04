package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/handlers"
	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/services/opengraphsvc"
	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func RootCmd() *cobra.Command {
	opts := &handlers.Options{
		Path:                "/v1",
		Port:                3000,
		ShutdownGracePeriod: 5 * time.Second,
	}
	deps := &handlers.Dependencies{
		Logger: logger.GetInstance(),
	}

	c := &cobra.Command{
		Use:   "api",
		Short: "serves the tenant REST API",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			// gormDB, _ := database.Connection()
			// deps.GormDB = gormDB
			openGraphSvc := opengraphsvc.Handler(deps.Logger)
			deps.Services.OpenGraphSvc = openGraphSvc

			service, serviceErr := handlers.NewService(ctx, opts, deps)
			if serviceErr != nil {
				return Cancel(serviceErr, cancel, service)
			}
			service.Start()
			deps.Logger.Info("api serving")
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				deps.Logger.Info("terminating: context canceled")
			case <-signals:
				deps.Logger.Info("terminating: via signal")
			}
			return Cancel(nil, cancel, service)
		},
	}

	return c
}

func Cancel(err error, cancel context.CancelFunc, closers ...io.Closer) error {
	if cancel != nil {
		cancel()
	}
	var eg errgroup.Group
	for i := range closers {
		closer := closers[i]
		if !isNil(closer) {
			eg.Go(closer.Close)
		}
	}
	waitErr := eg.Wait()
	if waitErr == nil {
		return err
	}
	if err == nil {
		return waitErr
	}
	return errors.Wrap(err, waitErr.Error())
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
