package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-router-tracker/config"
	"github.com/robertobadjio/tgtime-router-tracker/internal/background"
	kafkaModule "github.com/robertobadjio/tgtime-router-tracker/internal/kafka"
	timeService "github.com/robertobadjio/tgtime-router-tracker/internal/time"
	"github.com/robertobadjio/tgtime-router-tracker/internal/tracker"
)

var quit = make(chan struct{})
var checks = make(map[string]map[string]struct{})

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

	// msg="RPC error: context deadline exceeded, context deadline exceeded"
	//ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	//defer cancel()
	ctx := context.Background()
	bc.Start(ctx, logger)
	<-quit
}

func buildTrackerTaskFunc(
	timeClient *timeService.AggregatorClient,
	routerTracker *tracker.Tracker,
) func(ctx context.Context, logger log.Logger) {
	return func(ctx context.Context, logger log.Logger) {
		macAddresses, err := routerTracker.GetMacAddresses(ctx)
		if err != nil {
			_ = logger.Log("msg", err.Error())
			return
		}
		cfg := config.New()
		kafka := kafkaModule.NewKafka(cfg.KafkaHost + ":" + cfg.KafkaPort)

		currentDateTime := time.Now().Unix()

		// TODO: Нужна распределенная блокировка (через Redis) на случай, если роутеров будет n штук и n-инстансев текущего приложения
		// Чтобы один экземпляр приложения обрабатывал один роутер

		// TODO: Цикл по роутерам
		// TODO: Сохранять в БД и публиковать в Kafka бачами
		for _, macAddress := range macAddresses {
			fmt.Println("CYCLE", macAddress)
			// TODO: Пока роутер один
			// Если будет n роутеров, нужно ходить по каждому из них и забирать Mac-адреса активных устройств
			err = timeClient.CreateTime(ctx, macAddress, currentDateTime, 1)
			if err != nil {
				_ = logger.Log("msg", err.Error())
			}

			// TODO: Писать в Redis?
			dateNow := time.Now().Format("2006-01-02")
			_, ok := checks[dateNow]
			if !ok {
				checks[dateNow] = make(map[string]struct{})
			}

			_, ok = checks[dateNow][macAddress]
			if !ok {
				fmt.Println("PRODUCE", macAddress)
				err = kafka.ProduceInOffice(ctx, macAddress)
				if err != nil {
					_ = logger.Log("kafka", "produce in office message", "msg", err.Error())
				}
				checks[dateNow][macAddress] = struct{}{}
			}
			fmt.Println("PRODUCE", checks)
		}
	}
}
