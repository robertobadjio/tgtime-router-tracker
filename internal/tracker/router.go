package tracker

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gopkg.in/routeros.v2"
)

// Tracker Трекер роутера
type Tracker struct {
	host       string
	port       string
	login      string
	password   string
	properties string
}

// NewRouterTracker Конструктор трекер роутера
func NewRouterTracker(host, port, login, password string) *Tracker {
	properties := flag.String("properties", "mac-address", "Properties")

	return &Tracker{
		host:       host,
		port:       port,
		login:      login,
		password:   password,
		properties: *properties,
	}
}

// GetMacAddresses Получение списка mac-адресов подключенных к роутеру
func (r Tracker) GetMacAddresses(_ context.Context) ([]string, error) {
	var macAddresses []string
	flag.Parse()

	c, err := routeros.Dial(r.buildAddress(), r.login, r.password)
	if err != nil {
		log.Fatal(err) // TODO: !
	}
	defer c.Close()

	reply, err := c.Run("/interface/wireless/registration-table/print", "=.proplist="+r.properties)
	if err != nil {
		log.Fatal(err) // TODO: !
	}

	for _, re := range reply.Re {
		macAddresses = append(macAddresses, re.List[0].Value)
	}

	return macAddresses, nil
}

func (r Tracker) buildAddress() string {
	return fmt.Sprintf("%s:%s", r.host, r.port)
}
