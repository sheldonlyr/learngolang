package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	myqueue := queue.NewRingBuffer(4096)

	go func() {
		defer wg.Done()

		for index := 0; index < 10; index++ {
			ok, err := myqueue.Offer(index)
			if ok {
				time.Sleep(time.Second)
			} else {
				if err != nil {
					log.Println("Queue offer:", err)
					return
				}
				log.Println("Queue if full...")
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			index, err := myqueue.Get()
			if err != nil {
				log.Println("Queue get:", err)
				return
			}
			log.Println(index)
		}
	}()

	killSignal := <-interrupt
	fmt.Println()
	log.Println("Got signal:", killSignal)

	myqueue.Dispose()

	wg.Wait()
	log.Println("All goroutine are exit!")
}
