package tracker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
	modelRepo "github.com/robertobadjio/tgtime-router-tracker/internal/repository/model"
)

type routerService interface {
	GetMacAddresses() (map[uint][]string, error)
}

type kafka interface {
	ProduceInOffice(macAddress string) error
}

type aggregator interface {
	CreateTime(ctx context.Context, macAddress string, seconds, routerID int64) error
}

type routerRepo interface {
	GetAllActive(ctx context.Context) ([]modelRepo.Router, error)
}

// Tracker ...
type Tracker struct {
	routerService routerService
	kafka         kafka
	aggregator    aggregator
	routerRepo    routerRepo
	interval      time.Duration

	cancel context.CancelFunc

	mu     sync.Mutex
	checks map[string]map[string]struct{}
	ticker *time.Ticker
}

// NewTracker ...
func NewTracker(
	routerService routerService,
	kafka kafka,
	aggregator aggregator,
	routerRepo routerRepo,
	interval time.Duration,
) (*Tracker, error) {
	if routerService == nil {
		return nil, errors.New("routerService is nil")
	}

	if kafka == nil {
		return nil, errors.New("kafka is nil")
	}

	if aggregator == nil {
		return nil, errors.New("aggregator is nil")
	}

	if routerRepo == nil {
		return nil, errors.New("routerRepo is nil")
	}

	if interval <= 0 {
		return nil, errors.New("interval is invalid")
	}

	return &Tracker{
		routerService: routerService,
		kafka:         kafka,
		aggregator:    aggregator,
		routerRepo:    routerRepo,
		interval:      interval,
		checks:        make(map[string]map[string]struct{}),
	}, nil
}

// Run ...
func (t *Tracker) Run(ctx context.Context) error {
	if t.ticker != nil {
		return errors.New("tracker is already running")
	}

	t.mu.Lock()
	t.ticker = time.NewTicker(t.interval)
	t.mu.Unlock()

	defer t.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		select {
		case <-t.ticker.C:
			err := t.process(ctx)
			if err != nil {
				logger.Error(
					"component", "tracker",
					"during", "process",
					"error", err.Error(),
				)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Shutdown ...
func (t *Tracker) Shutdown() error {
	if t.cancel == nil {
		return fmt.Errorf("tracker not running")
	}

	t.cancel()

	return nil
}

func (t *Tracker) process(ctx context.Context) error {
	// TODO: Нужна распределенная блокировка (через Redis) на случай, если роутеров будет n штук и n-инстансев текущего приложения
	// Чтобы один экземпляр приложения обрабатывал один роутер

	currentDate := time.Now()

	routerMacAddresses, errGetMacAddresses := t.routerService.GetMacAddresses()
	if errGetMacAddresses != nil {
		logger.Error(
			"component", "tracker",
			"during", "get mac addresses",
			"err", errGetMacAddresses.Error(),
		)
	}

	for routerID, macAddresses := range routerMacAddresses {
		for _, macAddress := range macAddresses {
			// TODO: Batcher
			errCreateTime := t.aggregator.CreateTime(ctx, macAddress, currentDate.Unix(), int64(routerID)) // nolint : G115: integer overflow conversion uint -> int64 (gosec)
			if errCreateTime != nil {
				logger.Error(
					"component", "tracker",
					"during", "get mac addresses",
					"err", errCreateTime.Error(),
				)
			}

			if !t.setCheck(currentDate, macAddress) {
				errProduceInOffice := t.kafka.ProduceInOffice(macAddress)
				if errProduceInOffice != nil {
					logger.Error(
						"component", "kafka",
						"during", "produce in office message",
						"err", errProduceInOffice.Error(),
					)
				}
			}
		}
	}

	return nil
}

// TODO: Писать в Redis c TTL 24 часа
func (t *Tracker) setCheck(date time.Time, macAddress string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, found := t.checks[date.Format(time.DateOnly)]; !found {
		t.checks[date.Format(time.DateOnly)] = make(map[string]struct{})
	}

	_, found := t.checks[date.Format(time.DateOnly)][macAddress]
	if !found {
		t.checks[date.Format(time.DateOnly)][macAddress] = struct{}{}
	}

	return found
}
