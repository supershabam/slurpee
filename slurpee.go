package slurpee

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/url"
	"strings"
)

type slurpee struct {
	C    chan []byte
	Err  error
	conn redis.Conn
}

func NewSlurpee(redisUrl string, channel string) *slurpee {
	s := slurpee{}
	s.C = make(chan []byte)
	conn, err := newRedisConn(redisUrl)
	if err != nil {
		s.Err = err
		close(s.C)
		return &s
	}
	s.conn = conn
	psc := redis.PubSubConn{conn}
	psc.Subscribe(channel)
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				s.C <- v.Data
			case error:
				if v != io.EOF {
					s.Err = v
				}
				conn.Close()
				close(s.C)
				return
			}
		}
	}()
	return &s
}

func (s *slurpee) Stop() {
	s.conn.Close()
	return
}

func newRedisConn(uri string) (redis.Conn, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	var network string
	var host string
	var auth string

	switch u.Scheme {
	case "redis":
		network = "tcp"
		host = u.Host
		if !strings.Contains(host, ":") {
			host = fmt.Sprintf("%s:6379", host)
		}
		if u.User != nil {
			auth, _ = u.User.Password()
		}
	default:
		return nil, errors.New("invalid redis uri scheme")
	}

	conn, err := redis.Dial(network, host)
	if err != nil {
		return nil, err
	}

	if auth != "" {
		_, err = conn.Do("AUTH", auth)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
