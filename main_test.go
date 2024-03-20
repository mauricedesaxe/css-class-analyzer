package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	entry := PerformanceEntry{
		Timestamp:    time.Now(),
		Nanoseconds:  averageDuration.Nanoseconds(),
		Milliseconds: averageDuration.Milliseconds(),
		Seconds:      averageDuration.Seconds(),
	}

	if err := appendPerformanceResults(entry); err != nil {
		t.Fatalf("failed to append performance results: %s", err)
	}
}

type PerformanceEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Nanoseconds  int64     `json:"nanoseconds"`
	Milliseconds int64     `json:"milliseconds"`
	Seconds      float64   `json:"seconds"`
}

func appendPerformanceResults(entry PerformanceEntry) error {
	filePath := "./performance_results.json"
	var entries []PerformanceEntry

	// Read the existing entries
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		if err := json.Unmarshal(data, &entries); err != nil {
			return err
		}
	}

	// Append the new entry
	entries = append(entries, entry)

	// Write back to the file
	newData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, newData, 0644)
}
