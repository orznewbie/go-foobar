package time

import (
	"fmt"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 2)

	go func() {
		time.Sleep(3*time.Second)
		log.Info("Hello World!")
	}()

	select {
	case <-timer.C:
		log.Info("Timeout")
		return
	}
}

func TestAfterFunc(t *testing.T) {
	timer := time.AfterFunc(time.Second*3, printHello)
	time.Sleep(time.Second * 2)
	timer.Stop()
	var ch = make(chan struct{})
	<-ch
}

func printHello() {
	for i := 0; i < 10; i++ {
		fmt.Println("Hello ", i)
		time.Sleep(time.Second)
	}
}
