package daemon

import (
	"fmt"
	"testing"
	"time"
)

func daemontest() {
	Daemon(
		func() {
			for {
				fmt.Println("Daemon test.")
				time.Sleep(time.Second)
			}
		})
}

func noTest_daemon(t *testing.T) {
	daemontest()
}

func noTest_daemon_mt(t *testing.T) {
	for i := 0; i < 3; i++ {
		go daemontest()
	}
	time.Sleep(time.Hour * 1000)
}
