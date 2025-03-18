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
	"github.com/0x-Singularity/Augury/parser"
)

// QueryFakeula calls the Count FAKEula API and logs the query
func QueryFakeula(w http.ResponseWriter, r *http.Request) {
	// ✅ Fix CORS for React
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// ✅ Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Get IOC from query parameter
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	log.Println("Querying Count FAKEula for IOC:", ioc)

	// ✅ Get FAKEula API details from environment variables
	fakeulaURL := os.Getenv("FAKEULA_API_URL")
	fakeulaUser := os.Getenv("FAKEULA_USER")
	fakeulaPass := os.Getenv("FAKEULA_PASS")

	if fakeulaURL == "" {
		http.Error(w, "FAKEULA_API_URL is not set", http.StatusInternalServerError)
		return
	}

	// ✅ Build FAKEula request URL
	requestURL := fmt.Sprintf("%s%s", fakeulaURL, ioc)
	log.Println("FAKEula API Request URL:", requestURL)

	// ✅ Create HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(fakeulaUser, fakeulaPass) // Set authentication

	// ✅ Send request to FAKEula
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting FAKEula API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// ✅ Handle 404 error from FAKEula API
	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "FAKEula API returned 404 - Check API path", http.StatusNotFound)
		return
	}

	// ✅ Read response into buffer
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	// ✅ Decode JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(tee).Decode(&response); err != nil {
		log.Println("Warning: Unable to parse FAKEula response:", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	// ✅ Extract Result Count dynamically
	resultCount := parseResultCount(response)

	// ✅ Log the query in the database
	userName := r.Header.Get("X-User-Name") // Retrieve user from request headers
	if userName == "" {
		userName = "unknown"
	}

	if err := models.InsertQueryLog(ioc, resultCount, userName); err != nil {
		log.Println("Failed to log IOC lookup:", err)
	}

	// ✅ Retrieve the stored IOC log
	logEntry, err := models.GetQueryLog(ioc)
	if err != nil {
		log.Println("Failed to retrieve IOC log:", err)
	}

	// ✅ Append query log to the response
	if logEntry != nil {
		response["query_log"] = logEntry
	} else {
		response["query_log"] = "No log found"
	}

	// ✅ Pass the FAKEula response through the parser
	formattedResponse := parser.FormatFakeulaResponse(response)

	// ✅ Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	if err := json.NewEncoder(w).Encode(formattedResponse); err != nil {
		log.Println("Error encoding response:", err)
	}
}

// parseResultCount extracts the number of results dynamically
func parseResultCount(response map[string]interface{}) int {
	for _, value := range response {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			if arr, ok := value.([]interface{}); ok {
				return len(arr)
			}
		}
	}
	return 1 // Default to 1 if no array is found
}
