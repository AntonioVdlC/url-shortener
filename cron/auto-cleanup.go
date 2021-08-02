package cron

import (
	"log"
	"time"

	"url-shortener/db"
)

func AutoDeleteLinksJob() {
	db.DeleteOldLinks()

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			log.Println("Running 'AutoDeleteLinksJob'...")

			db.DeleteOldLinks()
		}
	}()
}
