package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"os"
	"tgtime-router-tracker/config"
	"tgtime-router-tracker/internal/background"
	timeService "tgtime-router-tracker/internal/time"
	"tgtime-router-tracker/internal/tracker"
	"time"
)

var quit = make(chan struct{})

func main() {
	cfg := config.New()

	var logger log.Logger

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	routerTracker := tracker.NewRouterTracker(
		cfg.RouterHost,
		cfg.RouterPort,
		cfg.RouterUserName,
		cfg.RouterPassword,
	)

	timeClient := timeService.NewTimeClient(*cfg, logger)

	bc := background.NewBackground(
		cfg.DelaySeconds,
		buildTrackerTaskFunc(timeClient, routerTracker),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bc.Start(ctx, logger)
	<-quit
}

func buildTrackerTaskFunc(
	timeClient *timeService.TimeClient,
	routerTracker *tracker.Tracker,
) func(ctx context.Context, logger log.Logger) {
	return func(ctx context.Context, logger log.Logger) {
		_ = logger.Log("msg", "Starting router tracker task...")
		macAddresses, err := routerTracker.GetMacAddresses(ctx)
		if err != nil {
			_ = logger.Log("msg", err.Error())
			return
		}

		currentDateTime := time.Now().Unix()
		// TODO: Цикл по роутерам
		for _, macAddress := range macAddresses {
			err = timeClient.CreateTime(ctx, macAddress, currentDateTime, 1) // TODO: Цикл по роутерам
			if err != nil {
				_ = logger.Log("msg", err.Error())
				return
			}
		}
	}
}
