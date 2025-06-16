package router_tracker

import (
	"flag"
	"fmt"

	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
	routerosInternal "github.com/robertobadjio/tgtime-router-tracker/internal/routeros"
)

// Tracker Трекер роутера.
type Tracker struct {
	clients    []routerosInternal.ClientInt
	properties string
	sentence   string
}

// NewRouterTracker Конструктор трекер роутера.
func NewRouterTracker(sentence string, routerClient []routerosInternal.ClientInt) (*Tracker, error) {
	properties := flag.String("properties", "mac-address", "Properties")
	flag.Parse() // TODO: Убрать?

	if sentence == "" {
		return nil, fmt.Errorf("invalid sentence")
	}

	if len(routerClient) == 0 {
		return nil, fmt.Errorf("invalid router clients")
	}

	return &Tracker{
		properties: *properties,
		sentence:   sentence,
		clients:    routerClient,
	}, nil
}

// GetMacAddresses Получение списка mac-адресов подключенных к роутеру.
func (t Tracker) GetMacAddresses() (map[uint][]string, error) {
	var macAddresses map[uint][]string

	for _, c := range t.clients {
		reply, errRun := c.Run(t.sentence, "=.proplist="+t.properties)
		if errRun != nil {
			logger.Error(
				"component", "router service",
				"during", "get mac addresses",
				"router ID", c.ID(),
				"err", errRun.Error(),
			)
			continue
		}

		macAddresses[c.ID()] = make([]string, 0, len(reply.Re))
		for _, re := range reply.Re {
			macAddresses[c.ID()] = append(macAddresses[c.ID()], re.List[0].Value)
		}
	}

	return macAddresses, nil
}
