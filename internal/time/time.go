package time

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	pb "github.com/robertobadjio/tgtime-aggregator/api/v1/pb/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tgtime-router-tracker/config"
)

type TimeClient struct {
	cfg    *config.Config
	logger log.Logger
}

func NewTimeClient(cfg config.Config, logger log.Logger) *TimeClient {
	return &TimeClient{cfg: &cfg, logger: logger}
}

func (tc TimeClient) CreateTime(ctx context.Context, macAddress string, seconds, routerId int64) error {
	client, err := grpc.NewClient(tc.buildAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}
	defer func() { _ = client.Close() }() // Игнорируем ошибку в явном виде

	timeAggregatorClient := pb.NewAggregatorClient(client)

	//ctxTemp, cancel := context.WithTimeout(ctx, 50*time.Second)
	//defer cancel()
	response, err := timeAggregatorClient.CreateTime(
		ctx,
		&pb.CreateTimeRequest{MacAddress: macAddress, Seconds: seconds, RouterId: routerId},
	)
	if err != nil {
		return fmt.Errorf("CreateTime: %v", err)
	}

	_ = tc.logger.Log("MacAddress: %s", response.MacAddress)

	return nil
}

func (tc TimeClient) buildAddress() string {
	return fmt.Sprintf("%s:%s", tc.cfg.TgTimeAggregatorHost, tc.cfg.TgTimeAggregatorPort)
}
