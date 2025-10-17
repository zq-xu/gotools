package utilsx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Notify
func Notify(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer func() {
			signal.Stop(c)
			close(c)
		}()

		sig := <-c
		fmt.Printf("exit by signal %v \n", sig)
		cancel()
	}()
}
