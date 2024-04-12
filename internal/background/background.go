package background

import (
	"context"
	"time"
)

type Background struct {
	delay time.Duration
	task  func(ctx context.Context)
}

func NewBackground(delay time.Duration, task func(ctx context.Context)) *Background {
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

func (b Background) Start(ctx context.Context) {
	ticker := time.NewTicker(b.delay)

	for {
		select {
		case <-ticker.C:
			go b.task(ctx)
		}
	}
}
