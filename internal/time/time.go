package time

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	"github.com/robertobadjio/tgtime-router-tracker/config"
)

// AggregatorClient Клиент gRPC-сервиса для сохранения времени
type AggregatorClient struct {
	cfg    *config.Config
	logger log.Logger
}

// NewTimeClient Конструктор gRPC-сервиса для сохранения времени
func NewTimeClient(cfg config.Config, logger log.Logger) *AggregatorClient {
	return &AggregatorClient{cfg: &cfg, logger: logger}
}

// CreateTime Сохранить времени "online" mac-адреса
func (tc AggregatorClient) CreateTime(
	ctx context.Context,
	macAddress string,
	seconds, routerID int64,
) error {
	client, err := grpc.NewClient(
		tc.buildAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}
	defer func() { _ = client.Close() }()

	timeAggregatorClient := pb.NewTimeV1Client(client)

	_, err = timeAggregatorClient.Create(
		ctx,
		&pb.CreateRequest{MacAddress: macAddress, Seconds: seconds, RouterId: routerID},
	)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			// Handle the error based on its status code
			if s.Code() == codes.NotFound {
				return fmt.Errorf("requested resource not found")
			}

			return fmt.Errorf("RPC error: %v", s.Message())
		}

		return fmt.Errorf("Non-RPC error: %v", err)
	}

	return nil
}

func (tc AggregatorClient) buildAddress() string {
	return fmt.Sprintf("%s:%s", tc.cfg.TgTimeAggregatorHost, tc.cfg.TgTimeAggregatorPort)
}
