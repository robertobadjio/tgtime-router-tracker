package app

import (
	"context"
	"github.com/robertobadjio/platform-common/pkg/closer"
	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	routerosInternal "github.com/robertobadjio/tgtime-router-tracker/internal/routeros"
	kafkaLib "github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/platform-common/pkg/db/pg"

	"github.com/robertobadjio/tgtime-router-tracker/internal/config"
	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
	routerRepo "github.com/robertobadjio/tgtime-router-tracker/internal/repository/router"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/kafka"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/router_tracker"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/time_aggregator"
	"github.com/robertobadjio/tgtime-router-tracker/internal/service/tracker"
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
		dbConfig, errNewDBConfig := config.NewDBConfig(sp.OS())
		if errNewDBConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init db config",
				"err", errNewDBConfig.Error(),
			)
		}

		dbClient := DBClient(ctx, dbConfig.DSN(), dbConfig.QueryTimeout())
		PGRouterRepo := routerRepo.NewPgRepository(dbClient)

		routers, errGetAllActive := PGRouterRepo.GetAllActive(ctx)
		if errGetAllActive != nil {
			logger.Fatal(
				"component", "di",
				"during", "get active routers",
				"err", errGetAllActive.Error(),
			)
		}

		routerTrackerConfig, errNewRouterTrackerConfig := config.NewRouterTrackerConfig(sp.OS())
		if errNewRouterTrackerConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init router tracker config",
				"err", errNewRouterTrackerConfig.Error(),
			)
		}

		routerOSClients := make([]routerosInternal.ClientInt, 0, len(routers))
		for _, r := range routers {
			routerOSClient, errRouterOSNewClient := routerosInternal.NewClient(
				r.Address,
				r.Login,
				r.Password,
				r.ID,
				routerTrackerConfig.DialTimeout(),
			)
			if errRouterOSNewClient != nil {
				logger.Error(
					"component", "di",
					"during", "create router OS client",
					"err", errRouterOSNewClient.Error(),
				)
				continue
			}

			closer.Add(func() error {
				routerOSClient.Close()
				return nil
			})

			routerOSClients = append(routerOSClients, routerOSClient)
		}

		routerService, errNewRouterTracker := router_tracker.NewRouterTracker(
			routerTrackerConfig.RegistrationTableSentence(),
			routerOSClients,
		)
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

		conn, errDialLeader := kafkaLib.DialLeader(
			ctx,
			"tcp",
			kafkaConfig.Address(),
			kafkaConfig.Topic(),
			kafkaConfig.Partition(),
		)
		if errDialLeader != nil {
			logger.Fatal(
				"component", "di",
				"during", "dial leader kafka",
				"err", errDialLeader.Error(),
			)
		}

		closer.Add(func() error { return conn.Close() })

		errSetWriteDeadline := conn.SetWriteDeadline(time.Now().Add(kafkaConfig.ConnDeadline()))
		if errSetWriteDeadline != nil {
			logger.Fatal(
				"component", "di",
				"during", "kafka set deadline",
				"err", errSetWriteDeadline.Error(),
			)
		}

		kf, errKafka := kafka.NewKafka(conn)
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
		client, errNewClient := grpc.NewClient(
			timeConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if errNewClient != nil {
			logger.Fatal(
				"component", "di",
				"during", "init time aggregator client",
				"err", errNewClient.Error(),
			)
		}
		closer.Add(func() error { return client.Close() })

		timeAggregatorClient := pb.NewTimeV1Client(client)
		ta, errNewTimeAggregator := time_aggregator.NewTimeAggregator(timeAggregatorClient)
		if errNewTimeAggregator != nil {
			logger.Fatal(
				"component", "di",
				"during", "init time aggregator service",
				"err", errNewTimeAggregator.Error(),
			)
		}

		cfgTracker, errNewTrackerConfig := config.NewTrackerConfig(sp.OS())
		if errNewTrackerConfig != nil {
			logger.Fatal(
				"component", "di",
				"during", "init tracker config",
				"err", errNewTrackerConfig.Error(),
			)
		}

		trackerService, errNewTracker := tracker.NewTracker(
			routerService,
			kf,
			ta,
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

// DBClient ...
func DBClient(ctx context.Context, connSrt string, queryTimeout time.Duration) db.Client {
	cl, err := pg.New(ctx, connSrt, queryTimeout)
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}

	err = cl.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}

	return cl
}
