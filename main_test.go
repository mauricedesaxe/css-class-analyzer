package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

func TestEndpointAccuracy(t *testing.T) {
	// Define the sample HTML input and expected classes
	sampleHTML := `<div class="bg-slate-950 text-gray-100 max-w-4xl mx-auto rounded-lg shadow-lg"><form class="flex flex-col gap-4"><div class="bg-red-500 p-2 rounded-md">Invalid HTML. Please enter valid HTML content.</div><textarea placeholder="Paste your HTML here" class="bg-slate-800 border border-slate-700 placeholder-gray-400 text-white rounded-md p-2">dialog</textarea><button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-500 disabled:cursor-not-allowed">Analyze</button></form><div class="flex flex-col gap-2 pt-6"></div></div>`
	expectedClasses := []string{"bg-slate-950", "text-gray-100", "max-w-4xl", "mx-auto", "rounded-lg", "shadow-lg", "flex", "flex-col", "gap-4", "bg-red-500", "p-2", "rounded-md", "bg-slate-800", "border", "border-slate-700", "placeholder-gray-400", "text-white", "bg-blue-500", "hover:bg-blue-700", "font-bold", "py-2", "px-4", "rounded", "disabled:bg-gray-500", "disabled:cursor-not-allowed", "gap-2", "pt-6"}

	// Start a mock fiber app with the postHTMLString handler
	app := fiber.New()
	app.Post("/", postHTMLString)

	// Prepare form data
	data := url.Values{}
	data.Set("html", sampleHTML)

	// Send the sample HTML to the / endpoint as form data
	req := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// the response body will be a JSON object with a `classNames` key containing an array of class names (strings)
	// parse the response body to get the class names
	var result struct {
		ClassNames []string `json:"classNames"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	receivedClasses := result.ClassNames
	log.Printf("Received classes: %v", receivedClasses)

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
