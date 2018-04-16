package main

import (
	"testing"
	"time"
)

func TestPingSlow(t *testing.T) {

	go PingSlow()
	go Ping()

	for i:=0; i < 10; i++ {
		PingReport()
		time.Sleep(time.Duration(1) *
			1000 * time.Millisecond)
	}

}