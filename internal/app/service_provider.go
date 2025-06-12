package app

import (
	"context"
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/platform-common/pkg/db/pg"
	"github.com/robertobadjio/tgtime-router-tracker/internal/config"
	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
	routerRepo "github.com/robertobadjio/tgtime-router-tracker/internal/repository/router"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/aggregator"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/kafka"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/router"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/tracker"
	"log"
	"time"
)

type serviceProvider struct {
	os config.OS

	tracker *tracker.Tracker
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// OS ...
func (sp *serviceProvider) OS() config.OS {
	if sp.os == nil {
		sp.os = config.NewOS()
	}

	return sp.os
}

// Tracker ...
func (sp *serviceProvider) Tracker(ctx context.Context) *tracker.Tracker {
	if sp.tracker == nil {
		cfgTracker, errNewTrackerConfig := config.NewTrackerConfig(sp.OS())
		if errNewTrackerConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init tracker config",
				"err", errNewTrackerConfig.Error(),
			)
		}

		routerService, errNewRouterTracker := router.NewRouterTracker()
		if errNewRouterTracker != nil {
			logger.Fatal(
				"component", "di",
				"during", "init router service",
				"err", errNewRouterTracker.Error(),
			)
		}

		kafkaConfig, errNewKafkaConfig := config.NewKafkaConfig(sp.OS())
		if errNewKafkaConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init kafka config",
				"err", errNewKafkaConfig.Error(),
			)
		}
		kf, errKafka := kafka.NewKafka(kafkaConfig.Address())
		if errKafka != nil {
			logger.Fatal(
				"component", "di",
				"during", "init kafka service",
				"err", errKafka.Error(),
			)
		}

		timeConfig, errNewAggregatorConfig := config.NewAggregatorConfig(sp.OS())
		if errNewAggregatorConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init aggregator config",
				"err", errNewAggregatorConfig.Error(),
			)
		}
		t := aggregator.NewAggregatorClient(timeConfig.Address())

		dbConfig, errNewDBConfig := config.NewDBConfig(sp.OS())
		if errNewDBConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init db config",
				"err", errNewDBConfig.Error(),
			)
		}

		dbClient := DBClient(ctx, dbConfig.DSN())
		PGRouterRepo := routerRepo.NewPgRepository(dbClient)

		trackerService, errNewTracker := tracker.NewTracker(
			routerService,
			kf,
			t,
			PGRouterRepo,
			cfgTracker.Interval(),
		)
		if errNewTracker != nil {
			logger.Fatal(
				"component", "di",
				"during", "init tracker service",
				"err", errNewTrackerConfig.Error(),
			)
		}

		sp.tracker = trackerService
	}

	return sp.tracker
}

func DBClient(ctx context.Context, connSrt string) db.Client {
	cl, err := pg.New(
		ctx,
		connSrt,
		1*time.Second,
	)
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}

	err = cl.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}

	return cl
}
