package routeros

import "gopkg.in/routeros.v2"

// ClientInt ...
type ClientInt interface {
	Runner
	Closer
}

// Runner ...
type Runner interface {
	Run(sentence ...string) (*routeros.Reply, error)
	ID() uint
}

// Closer ...
type Closer interface {
	Close()
}
