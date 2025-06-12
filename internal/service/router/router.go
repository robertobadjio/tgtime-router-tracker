package router

import (
	"context"
	"flag"
	"fmt"
	"time"

	"gopkg.in/routeros.v2"
)

const registrationTableSentence = "/interface/wireless/registration-table/print"
const dialTimeout = 10 * time.Second

// Tracker Трекер роутера.
type Tracker struct {
	properties string
}

// NewRouterTracker Конструктор трекер роутера.
func NewRouterTracker() (*Tracker, error) {
	properties := flag.String("properties", "mac-address", "Properties")
	flag.Parse()

	return &Tracker{
		properties: *properties,
	}, nil
}

// GetMacAddresses Получение списка mac-адресов подключенных к роутеру.
func (r Tracker) GetMacAddresses(_ context.Context, address, username, password string) ([]string, error) {
	fmt.Println("Prop list", r.properties)
	c, errDial := routeros.DialTimeout(address, username, password, dialTimeout)
	if errDial != nil {
		return []string{}, fmt.Errorf("dial: %w", errDial)
	}
	defer c.Close()

	reply, errRun := c.Run(registrationTableSentence, "=.proplist="+r.properties)
	if errRun != nil {
		return []string{}, fmt.Errorf("run: %w", errRun)
	}

	var macAddresses []string
	for _, re := range reply.Re {
		macAddresses = append(macAddresses, re.List[0].Value)
	}

	return macAddresses, nil
}
