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
	"github.com/0x-Singularity/Augury/parser"
)

// ExtractFromText receives a block of text, extracts IOCs, and queries FAKEula for each one
func ExtractFromText(w http.ResponseWriter, r *http.Request) {
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

	//--- Query Binary ---
	binaryURL := fmt.Sprintf("%scbr/binary/%s", baseURL, ioc)
	binaryData, err := fetchJSON(client, binaryURL, authUser, authPass)
	if err != nil {
		log.Printf("Binary query failed for %s: %v", ioc, err)
	} else {
		rawResponse["binary"] = binaryData
	}

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

// function to query all OIL sources
func QueryAllOIL(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%soil/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build OIL query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query OIL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var oilData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&oilData); err != nil {
		http.Error(w, "Error decoding OIL data", http.StatusInternalServerError)
		return
	}

	//Run oilData through the parser

	parsed := parser.FormatFakeulaResponse(oilData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// QueryPDNS queries the Passive DNS (PDNS) endpoint for a given IOC
func QueryPDNS(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%spdns/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build PDNS query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query PDNS", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var pdnsData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pdnsData); err != nil {
		http.Error(w, "Error decoding PDNS data", http.StatusInternalServerError)
		return
	}

	// Run PDNS data through the parser
	parsed := parser.FormatFakeulaResponse(pdnsData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// QueryLDAP queries the LDAP endpoint for a given IOC
func QueryLDAP(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%sldap/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build LDAP query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query LDAP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var ldapData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ldapData); err != nil {
		http.Error(w, "Error decoding LDAP data", http.StatusInternalServerError)
		return
	}

	// Run LDAP data through the parser
	parsed := parser.FormatFakeulaResponse(ldapData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// QueryGeoIP queries the GeoIP endpoint for a given IOC
func QueryGeoIP(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%sgeo/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build GeoIP query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query GeoIP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var geoData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		http.Error(w, "Error decoding GeoIP data", http.StatusInternalServerError)
		return
	}

	// Run GeoIP data through the parser
	parsed := parser.FormatFakeulaResponse(geoData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// QueryBinary queries the Binary endpoint for a given IOC
func QueryBinary(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%scbr/binary/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build Binary query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query Binary", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var binaryData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&binaryData); err != nil {
		http.Error(w, "Error decoding Binary data", http.StatusInternalServerError)
		return
	}

	// Run Binary data through the parser
	parsed := parser.FormatFakeulaResponse(binaryData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// QueryVPN queries the VPN endpoint for a given IOC
func QueryVPN(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%svpn/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build VPN query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query VPN", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var vpnData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&vpnData); err != nil {
		http.Error(w, "Error decoding VPN data", http.StatusInternalServerError)
		return
	}

	// No specific parser for VPN, return raw data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vpnData)
}

func QueryCBR(w http.ResponseWriter, r *http.Request) {
	ioc := r.URL.Query().Get("ioc")
	if ioc == "" {
		http.Error(w, "IOC parameter is required", http.StatusBadRequest)
		return
	}

	baseURL := os.Getenv("FAKEULA_API_URL")
	client := &http.Client{}
	user := os.Getenv("FAKEULA_USER")
	pass := os.Getenv("FAKEULA_PASS")

	url := fmt.Sprintf("%scbr/%s", baseURL, ioc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to build CBR query", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to query CBR", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cbrData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&cbrData); err != nil {
		http.Error(w, "Error decoding CBR data", http.StatusInternalServerError)
		return
	}

	//parse the CBR data
	log.Printf("CBR Data: %v", cbrData)
	parsed := parser.FormatFakeulaResponse(cbrData)
	log.Printf("Parsed CBR Data: %v", parsed)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}
