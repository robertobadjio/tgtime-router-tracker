package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/oklog/pkg/group"

	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
)

// App ...
type App struct {
	serviceProvider *serviceProvider
	gGroup          group.Group
}

// NewApp ...
func NewApp(ctx context.Context) (*App, error) {
	a := &App{
		gGroup: group.Group{},
	}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initCancelInterrupt,
		a.initTracker,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			logger.Fatal(
				"init", "deps",
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initTracker(ctx context.Context) error {
	tracker := a.serviceProvider.Tracker(ctx)
	a.gGroup.Add(func() error {
		logger.Info(
			"component", "app",
			"during", "run",
			"type", "tracker",
		)
		return tracker.Run(ctx)
	}, func(err error) {
		errShutdown := tracker.Shutdown()
		if errShutdown != nil {
			logger.Error(
				"component", "app",
				"during", "shutdown",
				"type", "tracker",
				"err", err.Error(),
			)
		}

		logger.Info(
			"component", "app",
			"during", "shutdown",
			"type", "tracker",
			"err", err.Error(),
		)
	})

	return nil
}

func (a *App) initCancelInterrupt(_ context.Context) error {
	cancelInterrupt := make(chan struct{})
	a.gGroup.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		logger.Info("component", "cancel interrupt", "err", "context canceled")
		close(cancelInterrupt)
	})

	return nil
}

// Run ...
func (a *App) Run() error {
	return a.gGroup.Run()
}
