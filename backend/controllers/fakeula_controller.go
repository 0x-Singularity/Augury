package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// FakeulaResponse represents the JSON response from FAKEula
type FakeulaResponse map[string]interface{}

// QueryFakeula calls the Count FAKEula API and returns results
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

	// Decode JSON response
	var response FakeulaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	// Return FAKEula response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
