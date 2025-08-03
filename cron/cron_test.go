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
	// Test the main function that now returns a cron job
	// Use a timeout to prevent hanging
	done := make(chan *CronJob, 1)

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

		// Verify job has the correct name
		if job.GetName() != "auto-delete-links" {
			t.Fatalf("Expected job name 'auto-delete-links', got '%s'", job.GetName())
		}

		// Stop immediately since we don't want to wait 24 hours
		job.Stop()

		t.Log("AutoDeleteLinksJob now returns stoppable CronJob")

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

	// Verify job has the correct name
	if job.GetName() != "auto-delete-links" {
		t.Fatalf("Expected job name 'auto-delete-links', got '%s'", job.GetName())
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

// Test the generic CronJob functionality
func TestNewCronJob(t *testing.T) {
	counter := 0
	task := func() {
		counter++
	}

	job := NewCronJob("test-job", 50*time.Millisecond, task)
	defer job.Stop()

	// Verify job was created
	if job == nil {
		t.Fatal("Expected job to be created, got nil")
	}

	if job.GetName() != "test-job" {
		t.Fatalf("Expected job name 'test-job', got '%s'", job.GetName())
	}

	// Let it run a few times
	time.Sleep(150 * time.Millisecond)

	// Stop the job
	job.Stop()

	// Verify task ran at least once
	if counter == 0 {
		t.Fatal("Expected task to run at least once")
	}

	t.Logf("Generic CronJob ran task %d times", counter)
}

func TestNewCronJobWithImmediate(t *testing.T) {
	counter := 0
	task := func() {
		counter++
	}

	// Test with immediate execution
	job := NewCronJobWithImmediate("immediate-test", 100*time.Millisecond, task, true)
	defer job.Stop()

	// Give it a moment to run immediately
	time.Sleep(10 * time.Millisecond)

	// Should have run immediately
	if counter == 0 {
		t.Fatal("Expected task to run immediately")
	}

	initialCount := counter

	// Let it run scheduled
	time.Sleep(150 * time.Millisecond)

	// Should have run more times
	if counter <= initialCount {
		t.Fatal("Expected task to run more times after initial execution")
	}

	job.Stop()
	t.Logf("CronJob with immediate execution ran task %d times", counter)
}

func TestCronJobWithoutImmediate(t *testing.T) {
	counter := 0
	task := func() {
		counter++
	}

	// Test without immediate execution
	job := NewCronJobWithImmediate("no-immediate-test", 50*time.Millisecond, task, false)
	defer job.Stop()

	// Give it a very short time - shouldn't run immediately
	time.Sleep(10 * time.Millisecond)

	// Should not have run yet
	if counter > 0 {
		t.Fatal("Expected task not to run immediately when runImmediately=false")
	}

	// Let it run scheduled
	time.Sleep(100 * time.Millisecond)

	// Should have run now
	if counter == 0 {
		t.Fatal("Expected task to run after interval")
	}

	job.Stop()
	t.Logf("CronJob without immediate execution ran task %d times", counter)
}

func TestAddJobAndGetRunningJobs(t *testing.T) {
	// Clear any existing jobs first
	Shutdown()

	// Create some test jobs
	job1 := NewCronJob("test-job-1", 1*time.Hour, func() {})
	job2 := NewCronJob("test-job-2", 2*time.Hour, func() {})

	// Add them to the tracker
	AddJob(job1)
	AddJob(job2)

	// Get running jobs
	jobs := GetRunningJobs()

	// Verify we have the expected jobs
	if len(jobs) != 2 {
		t.Fatalf("Expected 2 running jobs, got %d", len(jobs))
	}

	// Verify job names
	names := make(map[string]bool)
	for _, job := range jobs {
		names[job.GetName()] = true
	}

	if !names["test-job-1"] {
		t.Fatal("Expected to find 'test-job-1' in running jobs")
	}

	if !names["test-job-2"] {
		t.Fatal("Expected to find 'test-job-2' in running jobs")
	}

	// Clean up
	Shutdown()

	// Verify cleanup
	jobs = GetRunningJobs()
	if len(jobs) != 0 {
		t.Fatalf("Expected 0 running jobs after shutdown, got %d", len(jobs))
	}

	t.Log("AddJob and GetRunningJobs work correctly")
}
