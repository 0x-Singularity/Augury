package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/0x-Singularity/Augury/models"
)

// ExtractFromText receives a block of text, extracts IOCs, and queries FAKEula for each one
func ExtractFromText(w http.ResponseWriter, r *http.Request) {
	// CORS for React
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read input", http.StatusBadRequest)
		return
	}

	// Call FAKEula extractor
	extractURL := os.Getenv("FAKEULA_EXTRACT_URL")
	if extractURL == "" {
		extractURL = os.Getenv("FAKEULA_API_URL") + "extract"
	}
	req, err := http.NewRequest("POST", extractURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to create extract request", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(os.Getenv("FAKEULA_USER"), os.Getenv("FAKEULA_PASS"))
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call FAKEula extract", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var extractResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&extractResult); err != nil {
		http.Error(w, "Failed to parse extract response", http.StatusInternalServerError)
		return
	}

	iocs := extractIOCsFromResponse(extractResult)
	userName := r.Header.Get("X-User-Name")
	if userName == "" {
		userName = "unknown"
	}

	// Collect raw results before parsing
	rawResults := map[string]interface{}{}

	for _, ioc := range iocs {
		rawData, err := queryFakeulaForIOC(ioc, userName)
		if err != nil {
			log.Printf("Error processing IOC %s: %v", ioc, err)
			continue
		}
		rawResults[ioc] = rawData
	}

	// Run all raw IOC results through your parser
	//parsed := parser.FormatFakeulaResponse(rawResults)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": rawResults,
	})
}

func extractIOCsFromResponse(response map[string]interface{}) []string {
	iocs := []string{}
	data, ok := response["data"].([]interface{})
	if !ok {
		return iocs
	}
	for _, item := range data {
		if obj, ok := item.(map[string]interface{}); ok {
			if threat, ok := obj["threat"].(map[string]interface{}); ok {
				if indicator, ok := threat["indicator"].(map[string]interface{}); ok {
					if desc, ok := indicator["description"].(string); ok {
						iocs = append(iocs, desc)
					}
				}
			}
		}
	}
	return iocs
}

func queryFakeulaForIOC(ioc, userName string) (map[string]interface{}, error) {
	baseURL := os.Getenv("FAKEULA_API_URL")
	authUser := os.Getenv("FAKEULA_USER")
	authPass := os.Getenv("FAKEULA_PASS")
	client := &http.Client{}

	rawResponse := make(map[string]interface{})

	// --- Query Netflow ---
	netflowURL := fmt.Sprintf("%soil/netflow/%s", baseURL, ioc)
	netflowData, err := fetchJSON(client, netflowURL, authUser, authPass)
	if err != nil {
		log.Printf("Netflow query failed for %s: %v", ioc, err)
	} else {
		rawResponse["netflow"] = netflowData
	}

	// --- Query Security Logs (CoxSight) ---
	coxsightURL := fmt.Sprintf("%soil/coxsight/%s", baseURL, ioc)
	coxsightData, err := fetchJSON(client, coxsightURL, authUser, authPass)
	if err != nil {
		log.Printf("CoxSight query failed for %s: %v", ioc, err)
	} else {
		rawResponse["coxsight"] = coxsightData
	}

	// --- Query Asset Inventory ---
	assetURL := fmt.Sprintf("%sasset/%s", baseURL, ioc)
	assetData, err := fetchJSON(client, assetURL, authUser, authPass)
	if err != nil {
		log.Printf("Asset query failed for %s: %v", ioc, err)
	} else {
		rawResponse["asset"] = assetData
	}

	// --- Get PDNS Result Count ---
	resultCount := fetchPDNSResultCount(ioc)

	// --- Log the query ---
	if err := models.InsertQueryLog(ioc, resultCount, userName); err != nil {
		log.Println("Failed to log IOC lookup:", err)
	}

	// --- Retrieve and attach query logs ---
	logEntry, err := models.GetQueryLog(ioc)
	if err != nil {
		log.Println("Failed to retrieve IOC log:", err)
	}
	if logEntry != nil {
		var genericLogs []interface{}
		tmp, _ := json.Marshal(logEntry)
		_ = json.Unmarshal(tmp, &genericLogs)
		rawResponse["query_log"] = genericLogs
	} else {
		rawResponse["query_log"] = []interface{}{}
	}

	return rawResponse, nil
}

func fetchJSON(client *http.Client, url, user, pass string) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func fetchPDNSResultCount(ioc string) int {
	url := fmt.Sprintf("http://localhost:7000/pdns/%s/_summary", ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create PDNS summary request:", err)
		return 1
	}

	req.SetBasicAuth(os.Getenv("FAKEULA_USER"), os.Getenv("FAKEULA_PASS"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("PDNS summary request failed:", err)
		return 1
	}
	defer resp.Body.Close()

	var summary struct {
		Data []struct {
			NumResults int `json:"num_results"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		log.Println("Failed to decode PDNS summary:", err)
		return 1
	}

	if len(summary.Data) > 0 {
		log.Printf("PDNS Result Count for %s: %d", ioc, summary.Data[0].NumResults)
		return summary.Data[0].NumResults
	}

	log.Printf("No data found in PDNS summary for %s", ioc)
	return 1
}
