package main

import (
	"context"
	"fmt"
	"log"
	"tgtime-router-checker/config"
	"tgtime-router-checker/internal/background"
	timeService "tgtime-router-checker/internal/time"
	"tgtime-router-checker/internal/tracker"
	"time"
)

var quit = make(chan struct{})

func main() {
	cfg := config.New()
	routerTracker := tracker.NewRouterTracker(
		cfg.RouterHost,
		cfg.RouterPort,
		cfg.RouterUserName,
		cfg.RouterPassword,
	)

	timeClient := timeService.NewTimeClient(*cfg)

	bc := background.NewBackground(cfg.DelaySeconds, buildTrackerTaskFunc(timeClient, routerTracker))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bc.Start(ctx)
	<-quit
}

func buildTrackerTaskFunc(timeClient *timeService.TimeClient, routerTracker *tracker.Tracker) func(ctx context.Context) {
	return func(ctx context.Context) {
		fmt.Println("Starting router tracker task...")
		macAddresses, err := routerTracker.GetMacAddresses(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}

		currentDateTime := time.Now().Unix()
		// TODO: Цикл по роутерам
		for _, macAddress := range macAddresses {
			err = timeClient.CreateTime(ctx, macAddress, currentDateTime, 1) // TODO: Цикл по роутерам
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
