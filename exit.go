package exit_reload

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func ExitReload(prefix string, reload func(), breakdown func()) {
	done := make(chan struct{}, 1)
	sc := make(chan os.Signal, 1)

	//Protect from concurrent calling of handlers
	reloadingMutex := &sync.Mutex{}
	isRunning := true

	go func() {
		for {
			select {
			case a := <-sc:
				if a == syscall.SIGHUP {
					go func() {
						reloadingMutex.Lock()
						defer reloadingMutex.Unlock()

						//Perform a reload if still running
						if isRunning {
							fmt.Println()
							log.Printf("[%s] Reloading...\n", prefix)
							n := time.Now()

							reload()

							log.Printf("[%s] Took '%s' to reload\n", prefix, time.Now().Sub(n))
						}
					}()
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
	reloadingMutex.Lock()
	defer reloadingMutex.Unlock()
	isRunning = false
	fmt.Println()

	// Stop server
	log.Printf("[%s] Stopping...\n", prefix)
	n := time.Now()

	breakdown()

	log.Printf("[%s] Took '%s' to shutdown\n", prefix, time.Now().Sub(n))
	log.Printf("[%s] Goodbye\n", prefix)
}
