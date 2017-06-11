package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println("Hello world!")

	killChan := make(chan struct{})

	rp := &retryProcess{
		killChan:      killChan,
		shortInterval: 1 * time.Second,
		shortTotal:    10 * time.Second,
		longInterval:  3 * time.Second,
		longTotal:     30 * time.Second,
		msg:           "foo",
		fn: func() error {
			return fmt.Errorf("keep erroring")
		},
	}

	go func() {
		defer wg.Done()
		wg.Add(1)
		fmt.Println(rp.retry())
	}()

	<-time.After(3 * time.Second)
	fmt.Println("Close the killChan now")
	close(killChan)

	wg.Wait()
}

type retryProcess struct {
	killChan                  chan struct{}
	shortInterval, shortTotal time.Duration
	longInterval, longTotal   time.Duration
	msg                       string
	fn                        func() error
}

func (rp *retryProcess) retry() error {
	if err := rp.try(rp.shortInterval, rp.shortTotal); err != nil {
		fmt.Println("Switching to the long retry interval")
		return rp.try(rp.longInterval, rp.longTotal)
	}
	return nil
}

func (rp *retryProcess) try(interval, total time.Duration) error {
	var err = fmt.Errorf("process has not been executed yet")

	tickInterval := time.NewTicker(interval)
	tickTotal := time.NewTicker(total)
	defer tickTotal.Stop()
	defer tickInterval.Stop()
	fmt.Printf("Retrying every %s for a total of %s\n", interval.String(), total.String())

loop:
	for {
		select {
		case <-rp.killChan:
			fmt.Println("killChan closed")
			return fmt.Errorf("process was killed")
		case <-tickTotal.C:
			break loop
		case <-tickInterval.C:
			fmt.Println(rp.msg)
			if err = rp.fn(); err == nil {
				fmt.Println("Error is nil, breaking the loop")
				break loop
			}
		}
	}

	return err
}
