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

	// Open the log file to write to, wiping it clean on every start
	logFile, err := os.OpenFile("classes.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %s", err)
	}
	defer logFile.Close()

	// Create a new Scanner for the file
	scanner := bufio.NewScanner(file)

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
		classAttrValue := line[start+7 : start+7+end]

		// split the class names by space
		classNames := strings.Split(classAttrValue, " ")
		for _, className := range classNames {
			// write to log file
			if _, err := logFile.WriteString(className + "\n"); err != nil {
				log.Fatalf("failed to write to log file: %s", err)
			}
		}
	}

	// Check for errors during Scan. End of file is expected and not reported by Scan as an error.
	if err := scanner.Err(); err != nil {
		log.Fatalf("error during file scan: %s", err)
	}

}
