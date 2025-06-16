package time_aggregator

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
)

func TestTimeAggregator_NewTimeAggregator(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		timeAggregatorClient func() timeAggregatorClient

		expectedNilObj bool
		expectedErr    error
	}{
		"create config without time aggregator client": {
			timeAggregatorClient: func() timeAggregatorClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    ErrTimeAggregatorClientEmpty,
		},
		"create config": {
			timeAggregatorClient: func() timeAggregatorClient {
				timeAggregatorMock := NewMocktimeAggregatorClient(controller)
				require.NotNil(t, timeAggregatorMock)

				return timeAggregatorMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ta, err := NewTimeAggregator(test.timeAggregatorClient())
			require.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, ta)
			} else {
				assert.NotNil(t, ta)
			}
		})
	}
}

func TestTimeAggregator_CreateTime(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		timeAggregatorClient func() timeAggregatorClient

		expectedErr error
	}{
		"create time": {
			timeAggregatorClient: func() timeAggregatorClient {
				timeAggregatorMock := NewMocktimeAggregatorClient(controller)
				require.NotNil(t, timeAggregatorMock)

				timeAggregatorMock.
					EXPECT().
					Create(
						gomock.Any(),
						&pb.CreateRequest{MacAddress: "00-1A-2B-3C-4D-5E", Seconds: 1750087485, RouterId: 1},
					).
					Return(
						&pb.CreateResponse{MacAddress: "00-1A-2B-3C-4D-5E", Seconds: 1750087485, RouterId: 1},
						nil,
					)

				return timeAggregatorMock
			},

			expectedErr: nil,
		},
		"create time with non-RPC error": {
			timeAggregatorClient: func() timeAggregatorClient {
				timeAggregatorMock := NewMocktimeAggregatorClient(controller)
				require.NotNil(t, timeAggregatorMock)

				timeAggregatorMock.
					EXPECT().
					Create(
						gomock.Any(),
						&pb.CreateRequest{MacAddress: "00-1A-2B-3C-4D-5E", Seconds: 1750087485, RouterId: 1},
					).
					Return(
						nil,
						errors.New("grpc internal error"),
					)

				return timeAggregatorMock
			},

			expectedErr: fmt.Errorf("non-RPC error: %v", errors.New("grpc internal error")),
		},
		"create time with RPC error": {
			timeAggregatorClient: func() timeAggregatorClient {
				timeAggregatorMock := NewMocktimeAggregatorClient(controller)
				require.NotNil(t, timeAggregatorMock)

				timeAggregatorMock.
					EXPECT().
					Create(
						gomock.Any(),
						&pb.CreateRequest{MacAddress: "00-1A-2B-3C-4D-5E", Seconds: 1750087485, RouterId: 1},
					).
					Return(
						nil,
						status.Error(codes.ResourceExhausted, "resource has been exhausted"),
					)

				return timeAggregatorMock
			},

			expectedErr: fmt.Errorf("RPC error: %v", errors.New("resource has been exhausted")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ta, errNewTimeAggregator := NewTimeAggregator(test.timeAggregatorClient())
			require.NoError(t, errNewTimeAggregator)
			require.NotNil(t, ta)

			errCreateTime := ta.CreateTime(
				context.Background(),
				"00-1A-2B-3C-4D-5E",
				1750087485,
				1,
			)

			assert.Equal(t, test.expectedErr, errCreateTime)
		})
	}
}
