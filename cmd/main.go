package main

import (
	"fmt"
	"github.com/orznewbie/gotmpl/pkg/log"
	"sync"
	"sync/atomic"
)

var (
	dogCh  = make(chan struct{})
	fishCh = make(chan struct{})
	catCh  = make(chan struct{})
	wg     sync.WaitGroup
	num    uint32 = 2
)

func main() {
	wg.Add(3)
	go dog()
	go fish()
	go cat()
	catCh <- struct{}{}
	wg.Wait()
}

func dog() {
	var counter uint32
	for {
		if counter >= num {
			log.Info("dog done.")
			wg.Done()
			<-catCh
			return
		}
		<-catCh
		fmt.Println("dog")
		atomic.AddUint32(&counter, 1)
		dogCh <- struct{}{}
	}
}

func fish() {
	var counter uint32
	for {
		if counter >= num {
			log.Info("fish done.")
			wg.Done()
			return
		}
		<-dogCh
		fmt.Println("fish")
		atomic.AddUint32(&counter, 1)
		fishCh <- struct{}{}
	}
}

func cat() {
	var counter uint32
	for {
		if counter >= num {
			log.Info("cat done.")
			wg.Done()
			return
		}
		<-fishCh
		fmt.Println("cat")
		atomic.AddUint32(&counter, 1)
		catCh <- struct{}{}
	}
}
