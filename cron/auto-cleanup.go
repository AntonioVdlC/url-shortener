package cron

import (
	"log"
	"time"

	"url-shortener/db"
)

// runCleanupWithTimeout runs the cleanup with a timeout to prevent hanging
func runCleanupWithTimeout() {
	done := make(chan error, 1)

	go func() {
		done <- db.DeleteOldLinks()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error in cleanup: %v", err)
		}
	case <-time.After(10 * time.Second):
		log.Printf("Cleanup operation timed out after 10 seconds")
	}
}

// AutoDeleteLinksJobWithInterval creates a cleanup job with configurable interval
func AutoDeleteLinksJobWithInterval(interval time.Duration) *CronJob {
	return AutoDeleteLinksJobWithOptions(interval, true)
}

// AutoDeleteLinksJobWithOptions creates a cleanup job with configurable options
func AutoDeleteLinksJobWithOptions(interval time.Duration, runImmediately bool) *CronJob {
	return NewCronJobWithImmediate("auto-delete-links", interval, runCleanupWithTimeout, runImmediately)
}

// AutoDeleteLinksJob is the original function, now using the generic CronJob
func AutoDeleteLinksJob() *CronJob {
	return AutoDeleteLinksJobWithInterval(24 * time.Hour)
}
