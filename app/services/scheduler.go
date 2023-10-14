package services

import (
	"log"
	"time"
)

func Schedule(t string, callback func()) {
	for {
		if time.Now().Format("15:04") == t {
			log.Printf("Execute schedule [%s]\n", t)
			callback()
		}

		time.Sleep(1 * time.Minute)
	}
}
