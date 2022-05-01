package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 2)
	if timer == nil {

	}
	//time.Sleep(time.Second*3)
	//for {
	//	select {
	//	case <-timer.C:
	//		log.Info("Hello World!")
	//		return
	//	}
	//}
}

func TestAfterFunc(t *testing.T) {
	timer := time.AfterFunc(time.Second*3,printHello)
	time.Sleep(time.Second*2)
	timer.Stop()
	var ch = make(chan struct{})
	<- ch
}

func printHello() {
	for i := 0; i < 10; i++ {
		fmt.Println("Hello ", i)
		time.Sleep(time.Second)
	}
}
