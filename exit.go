package exit_reload

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ExitReload(prefix string, reload func(), breakdown func()) {
	done := make(chan struct{}, 1)
	sc := make(chan os.Signal, 1)

	go func() {
		for {
			select {
			case a := <-sc:
				if a == syscall.SIGHUP {
					go reload()
				} else {
					close(done)
					return
				}
			}
		}
	}()

	// Wait for exit signal
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-done
	fmt.Println()

	// Stop server
	log.Printf("[%s] Stopping...", prefix)
	n := time.Now()

	breakdown()

	log.Printf("[%s] Took '%s' to shutdown\n", prefix, time.Now().Sub(n))
	log.Printf("[%s] Goodbye\n", prefix)
}
