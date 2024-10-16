package background

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Background Сервис для фонового выполнения функции через определенные интервалы
type Background struct {
	delay time.Duration
	task  func(ctx context.Context, logger log.Logger)
}

// NewBackground Конструктор сервиса для фонового выполнения функций
func NewBackground(delay time.Duration, task func(ctx context.Context, logger log.Logger)) *Background {
	return &Background{delay: delay, task: task}
}

// Start Выполнить функцию
func (b Background) Start(ctx context.Context, logger log.Logger) {
	ticker := time.NewTicker(b.delay)

	for {
		select {
		case <-ticker.C:
			go b.task(ctx, logger)
		}
	}
}
