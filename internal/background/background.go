package background

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type Background struct {
	delay time.Duration
	task  func(ctx context.Context, logger log.Logger)
}

func NewBackground(delay time.Duration, task func(ctx context.Context, logger log.Logger)) *Background {
	return &Background{delay: delay, task: task}
}

/*func (b Background) Start() {
	ticker := time.NewTicker(b.delay)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				b.task()
			case <-stop:
				fmt.Println("Stopping Background")
				ticker.Stop()
				return
			}
		}
	}()

	return
}*/

func (b Background) Start(ctx context.Context, logger log.Logger) {
	ticker := time.NewTicker(b.delay)

	for {
		select {
		case <-ticker.C:
			go b.task(ctx, logger)
		}
	}
}
