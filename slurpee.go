package slurpee

import (
	"github.com/garyburd/redigo/redis"
	"io"
)

type slurpee struct {
	C    chan []byte
	Err  error
	conn redis.Conn
}

func NewSlurpee(redisUrl string, channel string) *slurpee {
	s := slurpee{}
	s.C = make(chan []byte)
	conn, err := redis.Dial("tcp", ":6379")
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
