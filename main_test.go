package main

import (
	"fmt"
	"testing"
	"time"
)

func TestMainSpeed(t *testing.T) {
	const runs = 100
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		startTime := time.Now()
		main()
		totalDuration += time.Since(startTime)
	}

	averageDuration := totalDuration / runs
	fmt.Printf("Average Duration over %d runs: %s\n", runs, averageDuration)
}
