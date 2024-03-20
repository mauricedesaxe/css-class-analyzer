package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()

	classesFromFile("index.html")

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Println("done in ", elapsed)
}

// read index.html file and serve each line to a new go routine
// each go routine will check if the line contains a `class` attribute
// if it does, it will classesFromFile the class names and log them (initially)
func classesFromFile(filename string) {
	// Open a file to read from
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var globalClassNames []string

	// Loop through all the lines
	for scanner.Scan() {
		// TODO we are assuming only HTML syntax.
		// TODO we are assuming only double quotes for attribute values.
		// TODO we are assuming only one class attribute per line.
		// TODO we are assuming the class attribute is not split across multiple lines.

		// Get the current line
		line := scanner.Text()

		// If the line does not contain a `class` attribute, skip it
		if !strings.Contains(line, "class=") {
			continue
		}

		// Find the start and end of the class attribute value in the line
		start := strings.Index(line, "class=\"")
		end := strings.Index(line[start+7:], "\"")

		// get the string between `class="` and `"` or `class='` and `'`
		// this is effectively the class names string
		classString := line[start+7 : start+7+end]

		// add the class names of this line to the class names of the whole file
		localClassNames := strings.Split(classString, " ")
		globalClassNames = append(globalClassNames, localClassNames...)
	}

	// Check for errors during Scan. End of file is expected and not reported by Scan as an error.
	if err := scanner.Err(); err != nil {
		log.Fatalf("error during file scan: %s", err)
	}

	// remove duplicates from the class names of the whole file using a map cause it's faster
	classMap := make(map[string]bool)
	var uniqueClassNames []string
	for _, className := range globalClassNames {
		if _, exists := classMap[className]; !exists {
			classMap[className] = true
			uniqueClassNames = append(uniqueClassNames, className)
		}
	}
	globalClassNames = uniqueClassNames

	// use buffered writing to log the class names in a freshly created (clean-wiped) file `classes.log`
	logFile, err := os.Create("classes.log")
	if err != nil {
		log.Fatalf("failed to create log file: %s", err)
	}
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)
	for _, className := range globalClassNames {
		_, err := writer.WriteString(className + "\n")
		if err != nil {
			log.Fatalf("failed to write to log file: %s", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("failed to flush writer: %s", err)
	}
}
