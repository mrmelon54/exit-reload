package exit_reload

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func c() chan struct{} {
	return make(chan struct{}, 1)
}

func TestExitReload(t *testing.T) {
	p, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)

	go func() {
		time.Sleep(5 * time.Second)
		t.Error("Failed tests")
		t.Fail()
		return
	}()

	reload, deconstruct := c(), c()

	go func() {
		time.Sleep(1 * time.Second)
		_ = p.Signal(syscall.SIGHUP)
		<-reload

		time.Sleep(1 * time.Second)
		_ = p.Signal(syscall.SIGINT)
		<-deconstruct
	}()

	ExitReload("TEST", func() {
		fmt.Println("Reload")
		close(reload)
	}, func() {
		fmt.Println("Deconstruct")
		close(deconstruct)
	})
}
