package cron

import "log"

var runningJobs []*CleanupJob

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
			job.Stop()
		}
	}
	runningJobs = nil
}
