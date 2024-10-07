package time

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"tgtime-router-tracker/config"
	"time"
)

type AggregatorClient struct {
	cfg    *config.Config
	logger log.Logger
}

func NewTimeClient(cfg config.Config, logger log.Logger) *AggregatorClient {
	return &AggregatorClient{cfg: &cfg, logger: logger}
}

func (tc AggregatorClient) CreateTime(ctx context.Context, macAddress string, seconds, routerId int64) error {
	client, err := grpc.NewClient(
		tc.buildAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithReturnConnectionError(),
		//grpc.WithKeepaliveParams(keepAlive),
		//grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}
	defer func() { _ = client.Close() }() // Игнорируем ошибку в явном виде

	timeAggregatorClient := pb.NewAggregatorClient(client)
	ctxTemp, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	_, err = timeAggregatorClient.CreateTime(
		ctxTemp,
		&pb.CreateTimeRequest{MacAddress: macAddress, Seconds: seconds, RouterId: routerId},
	)

	if err != nil {
		if status, ok := status.FromError(err); ok {
			// Handle the error based on its status code
			if status.Code() == codes.NotFound {
				return fmt.Errorf("Requested resource not found")
			} else {
				return fmt.Errorf("RPC error: %v, %v", status.Message(), ctxTemp.Err())
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
