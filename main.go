package main

import (
	"bufio"
	"css-class-analyzer/analyzer"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

func init() {
	// clean wipe the inputs and outputs directories
	err := os.RemoveAll("./inputs")
	if err != nil {
		fmt.Printf("Error removing files from inputs directory: %s\n", err)
	}
	err = os.RemoveAll("./outputs")
	if err != nil {
		fmt.Printf("Error removing files from outputs directory: %s\n", err)
	}
	if _, err := os.Stat("./inputs"); os.IsNotExist(err) {
		err := os.MkdirAll("./inputs", os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating inputs directory: %s\n", err)
		}
	}
	if _, err := os.Stat("./outputs"); os.IsNotExist(err) {
		err := os.MkdirAll("./outputs", os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating outputs directory: %s\n", err)
		}
	}
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// Get & sanitize the HTML input
		htmlInput := c.FormValue("html")
		if htmlInput == "" {
			return c.SendString("Please provide an HTML input")
		}
		p := bluemonday.UGCPolicy()
		p.AllowAttrs("class").Globally()
		sanitizedHTML := p.Sanitize(htmlInput)

		// Generate a unique request ID and timestamp
		requestId := uuid.New().String()
		timestamp := time.Now().Format("20060102-150405")

		// Create new directories for the input and output files
		inputDirName := fmt.Sprintf("./inputs/%s-%s", requestId, timestamp)
		err := os.MkdirAll(inputDirName, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error creating input directory: %s", err))
		}
		outputDirName := fmt.Sprintf("./outputs/%s-%s", requestId, timestamp)
		err = os.MkdirAll(outputDirName, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error creating output directory: %s", err))
		}

		// Write the sanitized HTML to a file
		htmlFileName := fmt.Sprintf("%s/input.html", inputDirName)
		file, err := os.Create(htmlFileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error creating file: %s", err))
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(sanitizedHTML)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error writing to file: %s", err))
		}
		err = writer.Flush()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error flushing writer: %s", err))
		}

		// Analyze the HTML input and return the log file
		logFileName := fmt.Sprintf("%s/classes.log", outputDirName)
		err = analyzer.Analyze(inputDirName, logFileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error analyzing HTML: %s", err))
		}

		// Start a new go routine that will delete the HTML and log files after a sleep duration
		go func() {
			time.Sleep(30 * time.Second)
			err := os.Remove(htmlFileName)
			if err != nil {
				fmt.Printf("Error deleting HTML file: %s\n", err)
			}
			err = os.Remove(logFileName)
			if err != nil {
				fmt.Printf("Error deleting log file: %s\n", err)
			}
		}()

		// Check file exists
		if _, err := os.Stat(logFileName); os.IsNotExist(err) {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error creating log file: %s", err))
		}

		// Read the log file and create a class list
		logFile, err := os.Open(logFileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error opening log file: %s", err))
		}
		defer logFile.Close()

		scanner := bufio.NewScanner(logFile)
		var classNames []string
		for scanner.Scan() {
			classNames = append(classNames, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error reading log file: %s", err))
		}

		// Return the log file
		return c.JSON(fiber.Map{
			"classNames": classNames,
		})
	})

	app.Listen(":3000")
}
