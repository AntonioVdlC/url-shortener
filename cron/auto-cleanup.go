package cron

import (
	"log"
	"time"

	"url-shortener/db"
)

// CleanupJob represents a running cleanup job that can be stopped
type CleanupJob struct {
	ticker *time.Ticker
	done   chan bool
}

// Stop gracefully stops the cleanup job
func (job *CleanupJob) Stop() {
	if job.ticker != nil {
		job.ticker.Stop()
	}
	if job.done != nil {
		job.done <- true
	}
}

// AutoDeleteLinksJobWithInterval creates a cleanup job with configurable interval
func AutoDeleteLinksJobWithInterval(interval time.Duration) *CleanupJob {
	return AutoDeleteLinksJobWithOptions(interval, true)
}

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

// AutoDeleteLinksJobWithOptions creates a cleanup job with configurable options
func AutoDeleteLinksJobWithOptions(interval time.Duration, runImmediately bool) *CleanupJob {
	// Set up recurring job
	ticker := time.NewTicker(interval)
	done := make(chan bool, 1) // Buffered to prevent blocking

	// Create job first
	job := &CleanupJob{
		ticker: ticker,
		done:   done,
	}

	go func() {
		// Run immediately if requested
		if runImmediately {
			select {
			case <-done:
				log.Println("Cleanup job stopped before initial run")
				return
			default:
				runCleanupWithTimeout()
			}
		}

		for {
			select {
			case <-ticker.C:
				log.Println("Running 'AutoDeleteLinksJob'...")
				runCleanupWithTimeout()
			case <-done:
				log.Println("Cleanup job stopped")
				return
			}
		}
	}()

	return job
}

// AutoDeleteLinksJob is the original function, now using the testable version
func AutoDeleteLinksJob() *CleanupJob {
	return AutoDeleteLinksJobWithInterval(24 * time.Hour)
}
