package analyzer

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	const runs = 100
	var totalDuration time.Duration
	var durations []time.Duration

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %s", err)
	}

	for i := 0; i < runs; i++ {
		startTime := time.Now()
		Analyze(cwd, "classes.log")
		totalDuration += time.Since(startTime)
		durations = append(durations, time.Since(startTime))
	}

	averageDuration := totalDuration / runs
	medianDuration := durations[runs/2]
	loc := loc()
	fileCount := fileCount()
	fmt.Printf("Did %d runs\n", runs)
	fmt.Printf("Total average Duration: %s\n", averageDuration)
	fmt.Printf("Total median Duration: %s\n", medianDuration)
	fmt.Printf("Average time per LoC: %v ns\n", averageDuration.Nanoseconds()/int64(loc))
	fmt.Printf("Average time per File: %v ns\n", averageDuration.Nanoseconds()/int64(fileCount))
}
