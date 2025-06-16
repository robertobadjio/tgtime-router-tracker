package time_aggregator

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
)

// ErrTimeAggregatorClientEmpty ...
var ErrTimeAggregatorClientEmpty = errors.New("time aggregator client must be set")

type timeAggregatorClient interface {
	Create(ctx context.Context, in *pb.CreateRequest, opts ...grpc.CallOption) (*pb.CreateResponse, error)
}

// TimeAggregator Клиент gRPC-сервиса для сохранения времени.
type TimeAggregator struct {
	timeAggregatorClient timeAggregatorClient
}

// NewTimeAggregator Конструктор gRPC-сервиса для сохранения времени.
func NewTimeAggregator(timeAggregatorClient timeAggregatorClient) (*TimeAggregator, error) {
	if timeAggregatorClient == nil {
		return nil, ErrTimeAggregatorClientEmpty
	}

	return &TimeAggregator{timeAggregatorClient: timeAggregatorClient}, nil
}

// CreateTime Сохранить времени "online" mac-адреса.
func (ta TimeAggregator) CreateTime(
	ctx context.Context,
	macAddress string,
	seconds, routerID int64,
) error {
	_, err := ta.timeAggregatorClient.Create(
		ctx,
		&pb.CreateRequest{MacAddress: macAddress, Seconds: seconds, RouterId: routerID},
	)

	if nil == err {
		return nil
	}

	if s, ok := status.FromError(err); ok {
		return fmt.Errorf("RPC error: %v", s.Message())
	}

	return fmt.Errorf("non-RPC error: %v", err)
}
