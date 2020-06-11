# gorillaradix

gorillaradix is a gorilla `sessions.Store` using [radix](https://github.com/mediocregopher/radix) as its persistence layer.

It supports vanilla redis and redis cluster

## Quickstart
To use this library you'd need to do something like:

### Vanilla redis
```go
store, err := gorillaradix.NewStore(gorillaradix.PoolConfiguration{
	Host: "my-redis-host:6379",
	ConnectionOptions: gorillaradix.ConnectionOptions{
		Password:    "terriblePassword",
		PingTimeout: 10 * time.Second,
		Timeout:     1 * time.Minute,
	    },
	}, gorillaradix.SessionOptions{
	Secret:      "terribleSecret",
})
````

### redis cluster
```go
store, err := gorillaradix.NewStoreCluster(gorillaradix.PoolConfiguration{
	Hosts: []string{"my-cluster-host-1:6379", "my-cluster-host-2:6379"},
	ConnectionOptions: gorillaradix.ConnectionOptions{
		Password:    "terriblePassword",
		PingTimeout: 10 * time.Second,
		Timeout:     1 * time.Minute,
	    },
	}, gorillaradix.SessionOptions{
	Secret:      "terribleSecret",
})
```

## Options
While there are reasonable defaults provided by the library, there are a number of variables to customize this session store.

### Notable session options:
| Variable | Description | Default |
| --- | --- | --- |
| `KeyPrefix` | Start session keys in redis with this string | `session_` |
| `MaxAge` | How long session data should be retained in seconds | `86400` |
| `MaxLength` | How much data should a session be able to store in bytes | `32768` |
| `Path` | What path should the cookie be set to | `/` |
| `Secret` | Secret (as a string) to validate session cookies | `notAGreatSecret` |


## Dependencies
* `github.com/gorilla/sessions`
* `github.com/gorilla/securecookie`
* `github.com/mediocregopher/radix`

## FAQ

### Why use the radix library?
It is one of the two client libraries [recommended by the redis project](https://redis.io/clients#go)
