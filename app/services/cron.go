package services

import (
	"log"
	"time"
)

func Cron(t string, callback func()) {
	for {
		if time.Now().Format("15:04") == t {
			log.Printf("Execute cron [%s]\n", t)
			callback()
		}

		time.Sleep(1 * time.Minute)
	}
}
