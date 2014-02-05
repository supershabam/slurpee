package main

import (
	"fmt"
	"github.com/supershabam/slurpee"
	"time"
)

func main() {
	s := slurpee.NewSlurpee("redis://localhost", "my-channel")
	go func() {
		select {
		case <-time.After(time.Second):
			s.Stop()
		}
	}()
	// s.C channel is closed when the redis connection dies, or slurpee is stopped
	for bytes := range s.C {
		fmt.Println("got bytes! %v", bytes)
	}
	// s.Err will be non-nil if there was a redis error
	if s.Err != nil {
		fmt.Printf("error: %v\n", s.Err)
	}
}
