package aggregator

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
)

// Aggregator Клиент gRPC-сервиса для сохранения времени.
type Aggregator struct {
	address string
}

// NewAggregatorClient Конструктор gRPC-сервиса для сохранения времени.
func NewAggregatorClient(address string) *Aggregator {
	return &Aggregator{address: address}
}

// CreateTime Сохранить времени "online" mac-адреса.
func (tc Aggregator) CreateTime(
	ctx context.Context,
	macAddress string,
	seconds, routerID int64,
) error {
	client, err := grpc.NewClient(
		tc.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("could not connect: %w", err)
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
