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
		err = Analyze(cwd, "classes.log")
		if err != nil {
			t.Fatalf("failed to analyze: %s", err)
		}
		totalDuration += time.Since(startTime)
		durations = append(durations, time.Since(startTime))
	}

	averageDuration := totalDuration / runs
	medianDuration := durations[runs/2]
	loc := loc(cwd)
	fileCount := fileCount(cwd)
	fmt.Printf("Did %d runs\n", runs)
	fmt.Printf("Total average Duration: %s\n", averageDuration)
	fmt.Printf("Total median Duration: %s\n", medianDuration)
	fmt.Printf("Average time per LoC: %v ns\n", averageDuration.Nanoseconds()/int64(loc))
	fmt.Printf("Average time per File: %v ns\n", averageDuration.Nanoseconds()/int64(fileCount))
}

func TestAccuracy(t *testing.T) {
	inputHTML := `<div class="bg-slate-950 text-gray-100 max-w-4xl mx-auto rounded-lg shadow-lg"><form class="flex flex-col gap-4"><div class="bg-red-500 p-2 rounded-md">Invalid HTML. Please enter valid HTML content.</div><textarea placeholder="Paste your HTML here" class="bg-slate-800 border border-slate-700 placeholder-gray-400 text-white rounded-md p-2">dialog</textarea><button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-500 disabled:cursor-not-allowed">Analyze</button></form><div class="flex flex-col gap-2 pt-6"></div></div>`
	expectedClasses := []string{"bg-slate-950", "text-gray-100", "max-w-4xl", "mx-auto", "rounded-lg", "shadow-lg", "flex", "gap-4", "bg-red-500", "p-2", "rounded-md", "bg-slate-800", "border", "border-slate-700", "placeholder-gray-400", "text-white", "p-2", "bg-blue-500", "hover:bg-blue-700", "text-white", "font-bold", "py-2", "px-4", "rounded", "disabled:bg-gray-500", "disabled:cursor-not-allowed", "flex", "gap-2", "pt-6"}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test.html")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the input HTML to the temporary file
	_, err = tmpFile.WriteString(inputHTML)
	if err != nil {
		t.Fatalf("failed to write to temp file: %s", err)
	}
	tmpFile.Close()

	// Extract classes from the temporary file
	classNames, err := classesFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to extract classes: %s", err)
	}

	// Compare the extracted classes with the expected classes
	// if there's a difference, print the difference
	for _, expectedClass := range expectedClasses {
		found := false
		for _, className := range classNames {
			if className == expectedClass {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected class %s not found", expectedClass)
		}
	}
}
