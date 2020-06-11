package gorillaradix

import (
	"github.com/gorilla/sessions"
	"github.com/mediocregopher/radix/v3"
	"time"
)

const (
	defaultMaxAge     = 86400             // how long should the session persist for
	defaultMaxLength  = 1024 * 32         // 32KB of session storage
	defaultKeyPrefix  = "session_"        // start session names with "session_"
	defaultCookiePath = "/"               // have it be site wide
	defaultSecret     = "notAGreatSecret" // default secret for session validation
)

type SessionOptions struct {
	sessions.Options
	MaxLength int    // Max amount of storage to be allowed in session data for a given user
	KeyPrefix string // within redis, what should session data keys be prefixed with
	Secret    string // Secret is for validating a session key
}

// Universal Redis Options
type ConnectionOptions struct {
	Password    string        // Password for the redis server
	PingTimeout time.Duration // How often should we ping the redis instance
	Timeout     time.Duration // Timeout on connecting
}

// helper function for creating the connection function for radix
// both cluster and regular redis will need this function
func createRadixConnectionFunction(configuration ConnectionOptions) (connectFunction radix.ConnFunc) {
	// If there is a password specified, we need to add it as part of the connection function
	if configuration.Password != "" {
		connectFunction = func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(configuration.Timeout),
				radix.DialAuthPass(configuration.Password),
			)
		}
	} else {
		connectFunction = func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(configuration.Timeout),
			)
		}
	}
	return
}

// Apply reasonable defaults to the options
func applyDefaultSessionConfiguration(options *SessionOptions) {
	if options.MaxLength == 0 {
		options.MaxLength = defaultMaxLength
	}
	if options.KeyPrefix == "" {
		options.KeyPrefix = defaultKeyPrefix
	}
	if options.Options.Path == "" {
		options.Options.Path = defaultCookiePath
	}
	if options.Options.MaxAge == 0 {
		options.Options.MaxAge = defaultMaxAge
	}
	if options.Secret == "" {
		options.Secret = defaultSecret
	}
}
