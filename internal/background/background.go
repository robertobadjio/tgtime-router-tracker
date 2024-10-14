package background

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Background struct {
	delay time.Duration
	task  func(ctx context.Context, logger log.Logger)
}

func NewBackground(delay time.Duration, task func(ctx context.Context, logger log.Logger)) *Background {
	return &Background{delay: delay, task: task}
}

func (b Background) Start(ctx context.Context, logger log.Logger) {
	ticker := time.NewTicker(b.delay)

	for {
		select {
		case <-ticker.C:
			go b.task(ctx, logger)
		}
	}
}
