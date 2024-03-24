package analyzer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func Analyze(dir string, output string) (err error) {
	startTime := time.Now()

	err = htmlFiles(dir, output)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(os.Args[0], ".test") {
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		loc := loc(dir)
		fileCount := fileCount(dir)

		fmt.Println("done in ", elapsed)
		fmt.Printf("time per loc: %v ns\n", elapsed.Nanoseconds()/int64(loc))
		fmt.Printf("time per file: %v ns\n", elapsed.Nanoseconds()/int64(fileCount))
	}
	return nil
}

// reads directory and children directories for html files and serves them to a function to get class names
// ultimately writes the class names to a log file
func htmlFiles(dir string, output string) (err error) {

	// use buffered writing to log the class names in a freshly created (clean-wiped) file
	logFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer logFile.Close()

	walkDirWg := sync.WaitGroup{}
	classStoreWg := sync.WaitGroup{}
	classNameChan := make(chan string, 1000) // Adjust buffer size as needed

	// walk the directory and serve each html file to the function
	// the function will return the class names of the file
	// the class names will be appended to a global slice
	var globalClassNames []string
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			walkDirWg.Add(1)
			go func(path string) {
				defer walkDirWg.Done()
				classNames, err := classesFromFile(path)
				if err != nil {
					log.Printf("error getting class names from file %q: %v\n", path, err)
				}
				for _, className := range classNames {
					classNameChan <- className // Send class names to the channel to be logged
				}
			}(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dir, err)
	}

	classStoreWg.Add(1)
	go func() {
		defer classStoreWg.Done()
		for className := range classNameChan {
			globalClassNames = append(globalClassNames, className)
		}
	}()
	walkDirWg.Wait()
	close(classNameChan)
	classStoreWg.Wait()

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

	// sort the class names
	slices.Sort(globalClassNames)

	// write the class names to the log file
	writer := bufio.NewWriter(logFile)
	for _, className := range globalClassNames {
		_, err := writer.WriteString(className + "\n")
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

// read index.html file and serve each line to a new go routine
// each go routine will check if the line contains a `class` attribute
// if it does, it will classesFromFile the class names and log them (initially)
func classesFromFile(filename string) (globalClassNames []string, err error) {
	// Open a file to read from
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the HTML file
	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
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

	return globalClassNames, nil
}

// gets lines of code for all files in dir/subdirs of ./pages
func loc(dir string) int {
	var lines int
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			lines += locFile(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dir, err)
	}
	return lines
}

// gets lines of code for a single file
func locFile(file string) int {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}
	lines := strings.Split(string(fileContent), "\n")
	return len(lines)
}

func fileCount(dir string) int {
	var count int
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			count++
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dir, err)
	}
	return count
}
