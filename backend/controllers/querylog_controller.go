package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-Singularity/Augury/models" // Import the model package
)

// LookupIOC handles IOC queries
func LookupIOC(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	// Check if IOC exists in database
	logEntry, err := models.GetQueryLog(ioc)
	if err != nil {
		http.Error(w, "Error retrieving query log", http.StatusInternalServerError)
		return
	}

	// If no entry found, return not found
	if logEntry == nil {
		http.Error(w, "IOC not found in logs", http.StatusNotFound)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

// LogIOC stores a new lookup entry
func LogIOC(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		IOC         string `json:"ioc"`
		ResultCount int    `json:"result_count"`
		UserName    string `json:"user_name"`
	}

	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure user_name is provided
	if requestData.UserName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	// Insert log into DB
	err = models.InsertQueryLog(requestData.IOC, requestData.ResultCount, requestData.UserName)
	if err != nil {
		http.Error(w, "Failed to log IOC lookup", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "IOC lookup logged successfully"})
}
