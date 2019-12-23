package grace

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Grace gives context and wait group for graceful shutdown
func Grace() (context.Context, *sync.WaitGroup) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()

	var wg sync.WaitGroup

	return ctx, &wg
}
