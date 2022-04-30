package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var (
		dogCh = make(chan struct{})
		fishCh = make(chan struct{})
		catCh = make(chan struct{})
		wg sync.WaitGroup
	)
	wg.Add(3)
	go dog(&wg, catCh, dogCh)
	go fish(&wg, dogCh, fishCh)
	go cat(&wg, fishCh, catCh)
	catCh <- struct{}{}
	wg.Wait()
}

func dog(wg *sync.WaitGroup, catCh, dogCh chan struct{}) {
	var counter int32
	for {
		if counter >= 100 {
			wg.Done()
			return
		}
		<- catCh
		fmt.Println("dog")
		atomic.AddInt32(&counter, 1)
		dogCh <- struct{}{}
	}
}

func fish(wg *sync.WaitGroup, dogCh, fishCh chan struct{}) {
	var counter int32
	for {
		if counter >= 100 {
			wg.Done()
			return
		}
		<- dogCh
		fmt.Println("fish")
		atomic.AddInt32(&counter, 1)
		fishCh <- struct{}{}
	}
}
	
func cat(wg *sync.WaitGroup, fishCh, catCh chan struct{}) {
	var counter int32
	for {
		if counter >= 100 {
			wg.Done()
			return
		}
		<- fishCh
		fmt.Println("cat")
		atomic.AddInt32(&counter, 1)
		catCh <- struct{}{}
	}
}