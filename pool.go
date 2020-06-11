package gorillaradix

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/mediocregopher/radix/v3"
)

// Vanilla Redis specific options
type PoolConfiguration struct {
	ConnectionOptions
	Host string
}

func NewStore(configuration PoolConfiguration, options SessionOptions) (sessions.Store, error) {
	connectFunction := createRadixConnectionFunction(configuration.ConnectionOptions)

	pool, err := radix.NewPool(
		"tcp",
		configuration.Host,
		10,
		radix.PoolConnFunc(connectFunction),
		radix.PoolPingInterval(configuration.PingTimeout),
	)

	applyDefaultSessionConfiguration(&options)

	return &RadixStore{
		redis:      pool,
		Codecs:     securecookie.CodecsFromPairs([]byte(options.Secret)),
		Options:    options,
		serializer: GobSerializer{},
	}, err
}
