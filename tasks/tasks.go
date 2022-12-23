package tasks

import (
	"log"
	"time"
)

func New(task func(stop chan bool), interval time.Duration) chan bool {
	log.Println("New task created")
	ticker := time.NewTicker(interval)
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				go task(stop)
			case <-stop:
				log.Println("Task exited")
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}
