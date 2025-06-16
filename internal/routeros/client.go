package routeros

import (
	"fmt"
	"time"

	routerosLib "gopkg.in/routeros.v2"
)

// Client ...
type Client struct {
	conn *routerosLib.Client
	id   uint
}

// NewClient ...
func NewClient(address, login, password string, id uint, dialTimeout time.Duration) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("invalid address")
	}

	if login == "" {
		return nil, fmt.Errorf("invalid login")
	}

	if password == "" {
		return nil, fmt.Errorf("invalid password")
	}

	if dialTimeout <= 0 {
		return nil, fmt.Errorf("invalid dial timeout")
	}

	if id <= 0 {
		return nil, fmt.Errorf("invalid ID router")
	}

	conn, errDialTimeout := routerosLib.DialTimeout(
		address,
		login,
		password,
		dialTimeout,
	)
	if errDialTimeout != nil {
		return nil, fmt.Errorf("dial router: %w", errDialTimeout)
	}

	return &Client{
		conn: conn,
		id:   id,
	}, nil
}

// Run ...
func (ros *Client) Run(sentence ...string) (*routerosLib.Reply, error) {
	return ros.Run(sentence...)
}

// ID ...
func (ros *Client) ID() uint {
	return ros.id
}

// Close ...
func (ros *Client) Close() {
	if ros.conn != nil {
		ros.conn.Close()
	}
}
