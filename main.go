package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
func htmlFiles(fn func(string) ([]string, error)) {
	cwd, _ := os.Getwd()
	log.Println("Current working directory:", cwd)

	err := filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			log.Printf("serving file: %s", path)
			fn(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", cwd, err)
	}
}

// read index.html file and serve each line to a new go routine
// each go routine will check if the line contains a `class` attribute
// if it does, it will classesFromFile the class names and log them (initially)
func classesFromFile(filename string) (globalClassNames []string, err error) {
	// Open a file to read from
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	// Parse the HTML file
	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML file: %s", err)
	}

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
		return nil, fmt.Errorf("failed to create log file: %s", err)
	}
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)
	for _, className := range globalClassNames {
		_, err := writer.WriteString(className + "\n")
		if err != nil {
			return nil, fmt.Errorf("failed to write to log file: %s", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return globalClassNames, nil
}
