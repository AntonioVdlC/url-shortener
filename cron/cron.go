package cron

import (
	"log"
	"time"
)

// CronJob represents a generic running cron job that can be stopped
type CronJob struct {
	name   string
	ticker *time.Ticker
	done   chan bool
}

// Stop gracefully stops the cron job
func (job *CronJob) Stop() {
	if job.ticker != nil {
		job.ticker.Stop()
	}
	if job.done != nil {
		job.done <- true
	}
}

// GetName returns the name of the cron job
func (job *CronJob) GetName() string {
	return job.name
}

// NewCronJob creates a new generic cron job
func NewCronJob(name string, interval time.Duration, task func()) *CronJob {
	ticker := time.NewTicker(interval)
	done := make(chan bool, 1) // Buffered to prevent blocking

	job := &CronJob{
		name:   name,
		ticker: ticker,
		done:   done,
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Running cron job '%s'...", name)
				task()
			case <-done:
				log.Printf("Cron job '%s' stopped", name)
				return
			}
		}
	}()

	return job
}

// NewCronJobWithImmediate creates a new cron job that optionally runs immediately
func NewCronJobWithImmediate(name string, interval time.Duration, task func(), runImmediately bool) *CronJob {
	ticker := time.NewTicker(interval)
	done := make(chan bool, 1) // Buffered to prevent blocking

	job := &CronJob{
		name:   name,
		ticker: ticker,
		done:   done,
	}

	go func() {
		// Run immediately if requested
		if runImmediately {
			select {
			case <-done:
				log.Printf("Cron job '%s' stopped before initial run", name)
				return
			default:
				log.Printf("Running initial execution of cron job '%s'...", name)
				task()
			}
		}

		for {
			select {
			case <-ticker.C:
				log.Printf("Running cron job '%s'...", name)
				task()
			case <-done:
				log.Printf("Cron job '%s' stopped", name)
				return
			}
		}
	}()

	return job
}

var runningJobs []*CronJob

func Init() {
	log.Println("Initialising cron jobs ...")

	job := AutoDeleteLinksJob()
	runningJobs = append(runningJobs, job)
}

// Shutdown gracefully stops all running cron jobs
func Shutdown() {
	log.Println("Shutting down cron jobs...")
	for _, job := range runningJobs {
		if job != nil {
			log.Printf("Stopping cron job '%s'", job.GetName())
			job.Stop()
		}
	}
	runningJobs = nil
}

// AddJob adds a cron job to the global job tracker
func AddJob(job *CronJob) {
	runningJobs = append(runningJobs, job)
	log.Printf("Added cron job '%s' to tracker", job.GetName())
}

// GetRunningJobs returns a copy of the currently running jobs
func GetRunningJobs() []*CronJob {
	jobs := make([]*CronJob, len(runningJobs))
	copy(jobs, runningJobs)
	return jobs
}
