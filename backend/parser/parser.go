package parser

import (
	"fmt"
)

// FakeulaResponse represents a parsed FAKEula response
type FakeulaResponse map[string]interface{}

// Cache to store query results
var resultsCache = make(map[string]FakeulaResponse)

// FormatFakeulaResponse parses and structures the FAKEula response to make it more readable and easier for the front end to display
func FormatFakeulaResponse(response FakeulaResponse) FakeulaResponse {
	formatted := make(FakeulaResponse)

	// Extract relevant fields from the FAKEula query
	// Checks if the key "data" exists in the response map and if its value can be cast into a slice ([]interface{})
	// If both conditions are true, the variable data will hold the value of response["data"] as a slice of empty interfaces (which means it can hold any type of value)
	if data, exists := response["data"].([]interface{}); exists {
		var parsedEntries []FakeulaResponse
		// Iterate through each element in the data slice
		// _, entry means that we ignore the index and just care about the entry value in each iteration
		for _, entry := range data {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				parsedEntry := FakeulaResponse{
					// Want to add an IOC field as well to organize queries with multiple IOCs in the input
					"callerIpAddress":   entryMap["callerIpAddress"],
					"coxAccountName":    entryMap["coxAccountName"],
					"displayName":       entryMap["displayName"],
					"oil":               entryMap["oil"],
					"timestamp":         entryMap["timestamp"],
					"userDisplayName":   entryMap["userDisplayName"],
					"userPrincipalName": entryMap["userPrincipalName"],
				}

				// I don't know how much any of this is necessary tbh, I don't like this at all and feel it's very inefficient, definitely want to change in the future
				// At first the results cache was not being filled with all the data returned in a FAKEula query, so I tried specifying fields for all the nested data,
				// however it just became a ton of if statements, which look terrible and still don't seem to be storing everything properly

				// Logic to handle nested client data
				if client, ok := entryMap["client"].(map[string]interface{}); ok {
					parsedEntry["client_as_org"] = client["as_org"]
					parsedEntry["client_asn"] = client["asn"]
					parsedEntry["client_ip"] = client["ip"]
				}
				// Handle binary hash response from /cbr/binary
				if binary, ok := entryMap["binary"].(map[string]interface{}); ok {
					parsedEntry["binary_md5"] = binary["md5"]
					parsedEntry["binary_sha256"] = binary["sha256"]
					parsedEntry["binary_filename"] = binary["filename"]
				}

				// Handle asset inventory
				if asset, ok := entryMap["asset"].(map[string]interface{}); ok {
					parsedEntry["asset_name"] = asset["name"]
					parsedEntry["asset_ip"] = asset["ip"]
					parsedEntry["asset_type"] = asset["type"]
				}

				// Handle GeoIP response
				if geo, ok := entryMap["geo"].(map[string]interface{}); ok {
					parsedEntry["geo_country"] = geo["country"]
					parsedEntry["geo_city"] = geo["city"]
					parsedEntry["geo_latitude"] = geo["latitude"]
					parsedEntry["geo_longitude"] = geo["longitude"]
				}

				// Handle LDAP response
				if ldap, ok := entryMap["ldap"].(map[string]interface{}); ok {
					parsedEntry["ldap_email"] = ldap["email"]
					parsedEntry["ldap_fullName"] = ldap["fullName"]
					parsedEntry["ldap_name"] = ldap["name"]
					parsedEntry["ldap_title"] = ldap["title"]
					parsedEntry["ldap_companyName"] = ldap["companyName"]
					parsedEntry["ldap_phone"] = ldap["phone"]
					parsedEntry["ldap_mobile"] = ldap["mobile"]
					parsedEntry["ldap_created"] = ldap["created"]
					parsedEntry["ldap_manager"] = ldap["manager"]
					parsedEntry["ldap_age"] = ldap["age"]
				}
				parsedEntries = append(parsedEntries, parsedEntry)
			}
		}
		formatted["formatted_data"] = parsedEntries
	} else {
		formatted["formatted_data"] = "No data available"
	}

	// Store formatted response in cache
	// %v converts formatted["formatted_data"] into a string regardless of type
	// I'm trying to use the original results as the key value like Garret suggested, not sure if I'm doing that properly
	cacheKey := fmt.Sprintf("%v", formatted["formatted_data"])
	resultsCache[cacheKey] = formatted

	// Temporary, prints the cache to the console to make sure the hash map is being occupied
	fmt.Println("=== Cached response added ===")
	PrintResultsCache()

	return formatted
}

// ParseFakeulaResponse calls FormatFakeulaResponse and returns the parsed data
/*func ParseFakeulaResponse(response FakeulaResponse) (FakeulaResponse, error) {
	// Check if length of response map is 0
	if len(response) == 0 {
		return nil, fmt.Errorf("empty FAKEula response")
	}
	return FormatFakeulaResponse(response), nil
}*/

func ParseFakeulaResponse(response FakeulaResponse) (FakeulaResponse, error) {
	// Check cache before processing
	cacheKey := fmt.Sprintf("%v", response["data"])
	if cachedResponse, found := resultsCache[cacheKey]; found {
		fmt.Println("Returning cached response")
		return cachedResponse, nil
	}

	// Process and store in cache
	if len(response) == 0 {
		return nil, fmt.Errorf("empty FAKEula response")
	}

	formatted := FormatFakeulaResponse(response)
	resultsCache[cacheKey] = formatted
	return formatted, nil
}

// Temp function to print the cache to see if it is being occupied
func PrintResultsCache() {
	fmt.Println("=== Printing resultsCache ===")
	if len(resultsCache) == 0 {
		fmt.Println("resultsCache is empty")
		return
	}

	for key, value := range resultsCache {
		fmt.Printf("Key: %s\nValue: %+v\n\n", key, value)
	}
	fmt.Println("=== End of resultsCache ===")
}
