package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Notify(cancel context.CancelFunc) {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Printf("exit by signal %v \n", s)
				cancel()
			}
		}
	}()
}
