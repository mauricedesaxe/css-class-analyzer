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
	sampleHTML := `<div class="bg-slate-950 text-gray-100 max-w-4xl mx-auto rounded-lg shadow-lg"><form class="flex flex-col gap-4"><div class="bg-red-500 p-2 rounded-md">Invalid HTML. Please enter valid HTML content.</div><textarea placeholder="Paste your HTML here" class="bg-slate-800 border border-slate-700 placeholder-gray-400 text-white rounded-md p-2">dialog</textarea><button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-500 disabled:cursor-not-allowed">Analyze</button></form><div class="flex flex-col gap-2 pt-6"></div></div>`
	expectedClasses := []string{"bg-slate-950", "text-gray-100", "max-w-4xl", "mx-auto", "rounded-lg", "shadow-lg", "flex", "flex-col", "gap-4", "bg-red-500", "p-2", "rounded-md", "bg-slate-800", "border", "border-slate-700", "placeholder-gray-400", "text-white", "bg-blue-500", "hover:bg-blue-700", "font-bold", "py-2", "px-4", "rounded", "disabled:bg-gray-500", "disabled:cursor-not-allowed", "gap-2", "pt-6"}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test.html")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the input HTML to the temporary file
	_, err = tmpFile.WriteString(sampleHTML)
	if err != nil {
		t.Fatalf("failed to write to temp file: %s", err)
	}
	tmpFile.Close()

	// Extract classes from the temporary file
	receivedClasses, err := classesFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to extract classes: %s", err)
	}

	// remove duplicates from the class names of the whole file using a map cause it's faster
	classMap := make(map[string]bool)
	var uniqueClassNames []string
	for _, className := range receivedClasses {
		if _, exists := classMap[className]; !exists {
			classMap[className] = true
			uniqueClassNames = append(uniqueClassNames, className)
		}
	}
	receivedClasses = uniqueClassNames

	// check that the length of `receivedClasses` is equal to the length of `expectedClasses`
	if len(receivedClasses) != len(expectedClasses) {
		t.Errorf("Expected %d classes, got %d", len(expectedClasses), len(receivedClasses))
	}

	// check that the classes in `expectedClasses` are present in the `receivedClasses`
	var expectedClassesThatWerentFound []string
	for _, className := range expectedClasses {
		if !contains(receivedClasses, className) {
			expectedClassesThatWerentFound = append(expectedClassesThatWerentFound, className)
		}
	}
	if len(expectedClassesThatWerentFound) > 0 {
		t.Errorf("Expected classes that weren't found: %s", expectedClassesThatWerentFound)
	}

	// check that there are no classes in `receivedClasses` that are not in `expectedClasses`
	var unexpectedClassesFound []string
	for _, className := range receivedClasses {
		if !contains(expectedClasses, className) {
			unexpectedClassesFound = append(unexpectedClassesFound, className)
		}
	}
	if len(unexpectedClassesFound) > 0 {
		t.Errorf("Unexpected classes found: %s", unexpectedClassesFound)
	}

	// log some information
	t.Logf("Expected classes that weren't found: %d", len(expectedClassesThatWerentFound))
	t.Logf("Unexpected classes found: %d", len(unexpectedClassesFound))
}

func contains(classes []string, class string) bool {
	for _, c := range classes {
		if c == class {
			return true
		}
	}
	return false
}
