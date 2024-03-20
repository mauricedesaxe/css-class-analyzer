package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	startTime := time.Now()

	htmlFiles(classesFromFile)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Println("done in ", elapsed)
}

// reads directory and children directories for html files and serves them to a function
func htmlFiles(fn func(string)) {
	// Open the current directory
	dir, err := os.Open(".")
	if err != nil {
		log.Fatalf("failed to open directory: %s", err)
	}
	defer dir.Close()

	// Read the directory
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalf("failed to read directory: %s", err)
	}

	// Serve the HTML files to the function
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			fn(file.Name())
		}
	}
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

	// Parse the HTML file
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("failed to parse HTML file: %s", err)
	}

	var globalClassNames []string

	// Define a recursive function to traverse the HTML nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					classNames := strings.Fields(a.Val)
					globalClassNames = append(globalClassNames, classNames...)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// Start traversing from the root node
	traverse(doc)

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
