package shared

import (
	"time"
)

func Ticker(tickTime time.Duration, stopTime time.Duration, onTick func()) {
	ticker := time.NewTicker(tickTime)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				onTick()
			}
		}
	}()

	time.Sleep(stopTime)
	ticker.Stop()
	done <- true
}
