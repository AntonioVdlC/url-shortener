package cron

import "log"

func Init() {
	log.Println("Initialising cron jobs ...")

	AutoDeleteLinksJob()
}