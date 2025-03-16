package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/0x-Singularity/Augury/models"
)

// FakeulaResponse represents the JSON response from FAKEula
type FakeulaResponse map[string]interface{}

// QueryFakeula calls the Count FAKEula API and logs the query
func QueryFakeula(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	log.Println("Querying Count FAKEula for IOC:", ioc)

	// Get FAKEula API details from environment variables
	fakeulaURL := os.Getenv("FAKEULA_API_URL")
	fakeulaUser := os.Getenv("FAKEULA_USER")
	fakeulaPass := os.Getenv("FAKEULA_PASS")

	if fakeulaURL == "" {
		http.Error(w, "FAKEULA_API_URL is not set", http.StatusInternalServerError)
		return
	}

	// Build FAKEula request URL
	requestURL := fmt.Sprintf("%s%s", fakeulaURL, ioc)
	log.Println("FAKEula API Request URL:", requestURL)

	// Create the HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Add basic authentication if required
	req.SetBasicAuth(fakeulaUser, fakeulaPass)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting FAKEula API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Ensure response is JSON, not an error page
	if resp.StatusCode == 404 {
		http.Error(w, "FAKEula API returned 404 - Check API path", http.StatusNotFound)
		return
	}

	// Read the response body into a buffer
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	// Decode JSON response for result count
	var response FakeulaResponse
	if err := json.NewDecoder(tee).Decode(&response); err != nil {
		log.Println("Warning: Unable to parse FAKEula response for result count:", err)
	}

	// Extract Result Count dynamically
	resultCount := parseResultCount(response)

	// Log the query in the database
	userName := r.Header.Get("X-User-Name") // Retrieve user from request headers
	if userName == "" {
		userName = "unknown"
	}

	err = models.InsertQueryLog(ioc, resultCount, userName)
	if err != nil {
		log.Println("Failed to log IOC lookup:", err)
	}

	// Retrieve the stored IOC log
	logEntry, err := models.GetQueryLog(ioc)
	if err != nil {
		log.Println("Failed to retrieve IOC log:", err)
	}

	// Append the logEntry to the response
	if logEntry != nil {
		response["query_log"] = logEntry
	} else {
		response["query_log"] = "No log found"
	}

	// Return updated JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	// Encode modified response with appended log data
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
	}
}

// parseResultCount attempts to extract the number of results dynamically
func parseResultCount(response FakeulaResponse) int {
	for _, value := range response {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			if arr, ok := value.([]interface{}); ok {
				return len(arr)
			}
		}
	}
	return 1 // Default to 1 if no array is found
}
