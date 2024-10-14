package time

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"github.com/robertobadjio/tgtime-router-tracker/config"
)

type AggregatorClient struct {
	cfg    *config.Config
	logger log.Logger
}

func NewTimeClient(cfg config.Config, logger log.Logger) *AggregatorClient {
	return &AggregatorClient{cfg: &cfg, logger: logger}
}

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

	timeAggregatorClient := pb.NewAggregatorClient(client)
	ctxTemp, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	_, err = timeAggregatorClient.CreateTime(
		ctxTemp,
		&pb.CreateTimeRequest{MacAddress: macAddress, Seconds: seconds, RouterId: routerID},
	)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			// Handle the error based on its status code
			if s.Code() == codes.NotFound {
				return fmt.Errorf("requested resource not found")
			} else {
				return fmt.Errorf("RPC error: %v, %v", s.Message(), ctxTemp.Err())
			}
		} else {
			// Handle non-RPC errors
			return fmt.Errorf("Non-RPC error: %v", err)
		}
	}

	return nil
}

func (tc AggregatorClient) buildAddress() string {
	return fmt.Sprintf("%s:%s", tc.cfg.TgTimeAggregatorHost, tc.cfg.TgTimeAggregatorPort)
}
