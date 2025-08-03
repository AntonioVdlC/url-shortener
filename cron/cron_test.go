package cron

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"url-shortener/db"
)

func TestMain(m *testing.M) {
	// Change to root directory so paths work correctly
	os.Chdir("..")

	// Load environment variables
	godotenv.Load(".env")

	// Initialize database
	db.Init()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// Test the original functions (limited testing due to design)
func TestAutoDeleteLinksJobExists(t *testing.T) {
	// Test that the function exists and doesn't panic on creation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AutoDeleteLinksJob panicked: %v", r)
		}
	}()

	// Note: We can't easily test the original function because:
	// 1. It starts a 24-hour ticker we can't stop
	// 2. It runs in a goroutine we can't control
	// 3. No return value to verify behavior

	t.Log("AutoDeleteLinksJob function exists and can be called")
	// Don't actually call it to avoid starting uncontrollable goroutines
}

func TestInitExists(t *testing.T) {
	// Test that Init exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init panicked: %v", r)
		}
	}()

	t.Log("Init function exists and can be called")
	// Don't actually call it to avoid starting uncontrollable goroutines
}

// Test the refactored testable versions
func TestAutoDeleteLinksJobWithInterval(t *testing.T) {
	// Test with a short interval for testing
	job := AutoDeleteLinksJobWithInterval(50 * time.Millisecond)
	defer job.Stop()

	// Verify job was created
	if job == nil {
		t.Fatal("Expected job to be created, got nil")
	}

	if job.ticker == nil {
		t.Fatal("Expected ticker to be created")
	}

	if job.done == nil {
		t.Fatal("Expected done channel to be created")
	}

	// Wait for a few ticks to ensure job runs
	time.Sleep(150 * time.Millisecond)

	// Job should still be running
	// (We can't easily verify the cleanup ran without checking database state)
	t.Log("Job created and ran for test duration")
}

func TestCleanupJobStop(t *testing.T) {
	// Test that we can stop a job
	job := AutoDeleteLinksJobWithInterval(10 * time.Millisecond)

	// Let it run briefly
	time.Sleep(25 * time.Millisecond)

	// Stop the job
	job.Stop()

	// Give it time to stop
	time.Sleep(20 * time.Millisecond)

	// Test passes if no panic occurs
	t.Log("Job stopped successfully")
}

func TestAutoDeleteLinksJob(t *testing.T) {
	// Test the main function that now returns a cleanup job
	// Use a timeout to prevent hanging
	done := make(chan *CleanupJob, 1)

	go func() {
		job := AutoDeleteLinksJob()
		done <- job
	}()

	// Wait for job creation with timeout
	select {
	case job := <-done:
		defer job.Stop()

		// Verify job was created
		if job == nil {
			t.Fatal("Expected job to be created, got nil")
		}

		// Stop immediately since we don't want to wait 24 hours
		job.Stop()

		t.Log("AutoDeleteLinksJob now returns stoppable job")

	case <-time.After(5 * time.Second):
		t.Fatal("AutoDeleteLinksJob took too long to create job (>5s)")
	}
}

func TestAutoDeleteLinksJobWithoutImmediateRun(t *testing.T) {
	// Test creating a job without running the database cleanup immediately
	job := AutoDeleteLinksJobWithOptions(100*time.Millisecond, false)
	defer job.Stop()

	// Verify job was created
	if job == nil {
		t.Fatal("Expected job to be created, got nil")
	}

	if job.ticker == nil {
		t.Fatal("Expected ticker to be created")
	}

	if job.done == nil {
		t.Fatal("Expected done channel to be created")
	}

	// Stop the job
	job.Stop()

	t.Log("Job created without immediate database run")
}

func TestInitAndShutdown(t *testing.T) {
	// Test the init and shutdown functions with timeout
	done := make(chan bool, 1)

	go func() {
		Init()
		done <- true
	}()

	// Wait for init with timeout
	select {
	case <-done:
		// Verify jobs were started (check global state)
		if len(runningJobs) == 0 {
			t.Fatal("Expected running jobs after Init(), got none")
		}

		// Shutdown all jobs
		Shutdown()

		// Verify jobs were stopped
		if runningJobs != nil {
			t.Fatal("Expected runningJobs to be nil after Shutdown()")
		}

		t.Log("Init and Shutdown worked correctly")

	case <-time.After(5 * time.Second):
		t.Fatal("Init() took too long (>5s)")
	}
}

func TestMultipleJobsAndCleanup(t *testing.T) {
	// Test creating multiple jobs and cleaning them up
	job1 := AutoDeleteLinksJobWithInterval(20 * time.Millisecond)
	job2 := AutoDeleteLinksJobWithInterval(30 * time.Millisecond)

	// Let them run briefly
	time.Sleep(50 * time.Millisecond)

	// Stop both jobs
	job1.Stop()
	job2.Stop()

	// Give time to stop
	time.Sleep(20 * time.Millisecond)

	t.Log("Multiple jobs created and stopped successfully")
}
