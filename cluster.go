package gorillaradix

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/mediocregopher/radix/v3"
)

// Redis Cluster specific options
type ClusterConfiguration struct {
	ConnectionOptions
	Hosts []string
}

func NewStoreCluster(configuration ClusterConfiguration, options SessionOptions) (sessions.Store, error) {
	connectFunction := createRadixConnectionFunction(configuration.ConnectionOptions)

	poolFunc := func(network, addr string) (radix.Client, error) {
		return radix.NewPool(
			network,
			addr,
			100,
			radix.PoolConnFunc(connectFunction),
			radix.PoolPingInterval(configuration.PingTimeout))
	}

	pool, err := radix.NewCluster(configuration.Hosts, radix.ClusterPoolFunc(poolFunc))
	applyDefaultSessionConfiguration(&options)

	return &RadixStore{
		redis:      pool,
		Codecs:     securecookie.CodecsFromPairs([]byte(options.Secret)),
		Options:    options,
		serializer: GobSerializer{},
	}, err
}
