package gorillaradix

import (
	"encoding/base32"
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/mediocregopher/radix/v3"
	"net/http"
	"strconv"
	"strings"
)

// Amount of time for cookies/redis keys to expire.
const (
	ErrSessionDataTooBig = "gorillaradix.RadixStore.save: data to persist is over size limit"
)

type RadixStore struct {
	redis      radix.Client
	Codecs     []securecookie.Codec
	Options    SessionOptions
	serializer SessionSerializer
}

// Get should return a cached session.
func (rs *RadixStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(rs, name)
}

// New creates new session and if possible, loads the stored data into it
func (rs *RadixStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var err error
	var ok bool

	// Create a new session and load session options into it
	session := sessions.NewSession(rs, name)
	session.Options = &rs.Options.Options
	// Mark the session as new for now
	session.IsNew = true

	// try to load a cookie, if so continue
	if c, errCookie := r.Cookie(name); errCookie == nil {
		// cookie was retrieved, let's try to decode the session id
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, rs.Codecs...)
		if err == nil {
			// Let's try to load the session data from the redis store
			ok, err = rs.load(session)
			// if ok is true and there is no error, we were successful
			if err == nil && ok {
				session.IsNew = false
			}
		}
	}
	return session, err

}

// Save should persist session to redis
func (rs *RadixStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	// if max age is greater than zero, we want to persist the data to redis
	if s.Options.MaxAge > 0 {
		// Check if an id has not yet been given
		if s.ID == "" {
			// Generates random id
			s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}

		// Let's try persisting the data to redis
		if err := rs.save(s); err != nil {
			return err
		}

		// Now let's generate the contents to persist as cookie to the browser
		encoded, err := securecookie.EncodeMulti(s.Name(), s.ID, rs.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(s.Name(), encoded, s.Options))
	} else {
		// So we need to delete the data from redis and browser
		// First let's delete the key from redis
		if err := rs.delete(s); err != nil {
			return err
		}
		// Now let's remove the session name from the browser
		http.SetCookie(w, sessions.NewCookie(s.Name(), "", s.Options))
	}
	return nil
}

// helper functions

// Generates a session key name using the session key prefix and the session id
func (rs *RadixStore) generateSessionKeyName(id string) string {
	return rs.Options.KeyPrefix + id
}

// Stores the session in redis.
func (rs *RadixStore) save(session *sessions.Session) error {
	serializedBytes, err := rs.serializer.Serialize(session)
	if err != nil {
		return err
	}

	// Ensuring that we're not storing too much data
	if rs.Options.MaxLength != 0 && len(serializedBytes) > rs.Options.MaxLength {
		return errors.New(ErrSessionDataTooBig)
	}

	// Will generate a redis command looking like
	//    SETEX session_sessionidvalue 86400 string-escaped-serialized-data
	// If we're storing our data for one day (86400 seconds)
	return rs.redis.Do(radix.Cmd(
		nil,
		"SETEX",
		rs.generateSessionKeyName(session.ID),
		strconv.Itoa(session.Options.MaxAge),
		string(serializedBytes),
	))
}

// Load reads the session from redis.
// returns true if there is session data in DB
func (rs *RadixStore) load(session *sessions.Session) (bool, error) {
	var data []byte
	// Will generate a redis command looking like
	//    GET session_sessionidvalue
	err := rs.redis.Do(radix.Cmd(
		&data,
		"GET",
		rs.generateSessionKeyName(session.ID),
	))

	// If there's no data or we get an error, let's return false and the error (if any)
	if err != nil || data == nil {
		return false, err
	}

	// Got data, deserialize it and return any deserialization errors
	return true, rs.serializer.Deserialize(data, session)
}

// Delete removes a session's contents from redis
func (rs *RadixStore) delete(session *sessions.Session) error {
	// Will generate a redis command looking like
	//    DEL session_sessionidvalue
	return rs.redis.Do(radix.Cmd(
		nil,
		"DEL",
		rs.generateSessionKeyName(session.ID),
	))
}
