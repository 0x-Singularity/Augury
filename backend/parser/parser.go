package parser

import (
	"encoding/json"
	"fmt"
)

// MultiLevelMap is the data structure to store parsed FAKEula data.

// - First level key: source
// - Second level key: structure type
// - Value: Slice of FakeulaEntry structs containing the actual data
type MultiLevelMap map[string]map[string][]FakeulaEntry

// ResultsCache defines a map type that stores parsed FAKEula responses for reuse
// Key is the stringified JSON data, Value is the parsed MultiLevelMap
type ResultsCache map[string]MultiLevelMap

//---------------------------Structs to represent different endpoint results from a FAKEula query-------------------------------------------------------------

// FakeulaEntry represents a parsed Fakeula response entry.
// This is the main struct that contains all the different types of data that can be returned from a FAKEula query
// Most fields are pointers so they can be nil if not present.
type FakeulaEntry struct {
	Oil    *OilInfo    `json:"oil"`
	Client *ClientInfo `json:"client,omitempty"`
	Binary *BinaryInfo `json:"binary,omitempty"`
	Asset  *AssetInfo  `json:"asset,omitempty"`
	Geo    *GeoInfo    `json:"geo,omitempty"`
	LDAP   *LdapInfo   `json:"ldap,omitempty"`
	PDNS   *PDNSInfo   `json:"pdns,omitempty"`
}

// OilInfo
type OilInfo struct {
	Timestamp     string `json:"timestamp"`
	UserPrincipal string `json:"userPrincipalName"`
	DisplayName   string `json:"displayName"`
	ClientIP      string `json:"clientIp"`
	ClientASNOrg  string `json:"clientAsOrg"`
	EventType     string `json:"eventType"`
	Outcome       string `json:"outcome"`
	Message       string `json:"message"`
}

// ClientInfo represents network client information
type ClientInfo struct {
	AsOrg string `json:"as_org"`
	ASN   int    `json:"asn"`
	IP    string `json:"ip"`
}

// BinaryInfo struct to match the CBR JSON structure, this one has a lot of nested stuff
type BinaryInfo struct {
	MD5        string   `json:"md5"`
	SHA256     string   `json:"sha256"`
	Filename   string   `json:"filename"`
	Accessed   string   `json:"accessed"`
	Hosts      []string `json:"hosts"`
	CodeSigned bool     `json:"codeSigned"`
	URL        string   `json:"url"`
}

// They said they don't like their current method of asset inventory, we may want to try and expand on how we present the data
type AssetInfo struct {
	Name          string `json:"name"`
	IP            string `json:"ip"`
	PlatformName  string `json:"platformName"`
	PlatformOwner string `json:"platformOwner"`
	Executive     string `json:"executive"`
	StackName     string `json:"stackName"`
	StackOwner    string `json:"stackOwner"`
	Created       string `json:"created"`
	Updated       string `json:"updated"`
}

type GeoInfo struct {
	CountryCode string  `json:"countryCode"`
	CountryName string  `json:"countryName"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ASNumber    string  `json:"asNumber"`
	ASOrg       string  `json:"asOrg"`
}

type LdapInfo struct {
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	CompanyName string `json:"companyName"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	Created     string `json:"created"`
	Manager     string `json:"manager"`
	Age         string `json:"age"`
}

// DNSAnswer represents a single DNS record answer
type DNSAnswer struct {
	Data  string `json:"data"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Count int    `json:"count"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// PDNSInfo contains Passive DNS information (historical DNS records)
type PDNSInfo struct {
	Answers []DNSAnswer `json:"answers"`
}

// result structure
type ParsedFakeulaResult struct {
	Data MultiLevelMap `json:"data"`
}

//--------------------Functions to parse and format the FAKEula response---------------------------------------------------------------------

// Global cache variable to store parsed results to avoid re-parsing
var resultsCache = make(ResultsCache)

// FormatFakeulaResponse parses and organizes the FAKEula response

func FormatFakeulaResponse(response map[string]interface{}) ParsedFakeulaResult {
	var parsedData = make(MultiLevelMap)

	// Check if "data" field exists in response
	if data, exists := response["data"].([]interface{}); exists {
		// Iterate through each entry in the data array
		for _, entry := range data {
			// Try to convert the entry to a map
			if entryMap, ok := entry.(map[string]interface{}); ok {
				// Create a FakeulaEntry struct and populate it with data from the entry map
				parsedEntry := FakeulaEntry{
					Oil:    parseOil(entryMap),
					Client: parseClient(entryMap),
					Binary: parseBinary(entryMap),
					Asset:  parseAsset(entryMap),
					Geo:    parseGeo(entryMap),
					LDAP:   parseLdap(entryMap),
					PDNS:   parsePDNS(entryMap),
				}

				// Extract keys for organizing the data in the MultiLevelMap
				//ioc := getString(entryMap, "key") //get ioc from the key since that is the queried IOC
				//oil := parsedEntry.Oil
				source := getSource(entryMap)

				structureTypes := getStructureTypes(&parsedEntry)

				for _, structType := range structureTypes {

					// Initialize nested maps if they don't exist
					// Level 1 - source
					if _, exists := parsedData[source]; !exists {
						parsedData[source] = make(map[string][]FakeulaEntry)
					}
					// Level 2 - structure type
					if _, exists := parsedData[source][structType]; !exists {
						parsedData[source][structType] = []FakeulaEntry{}
					}

					// Append the parsed entry to the appropriate slice in the MultiLevelMap
					parsedData[source][source] = append(parsedData[source][structType], parsedEntry)
				}
			}
		}
	}

	// Store the parsed data in the cache using the original data as the key
	cacheKey, err := json.Marshal(response["data"])
	if err == nil {
		resultsCache[string(cacheKey)] = parsedData
	}
	// Print cache for debugging
	//fmt.Println("=== Cached response added ===")
	//PrintResultsCache()

	return ParsedFakeulaResult{
		Data: parsedData,
	}

}

// Function to return a slice of structure types present in a FAKEula entry, this function determines the second level of the multi level map
func getStructureTypes(entry *FakeulaEntry) []string {
	structTypes := []string{}

	if entry.Oil != nil {
		structTypes = append(structTypes, "oil")
	}
	if entry.Client != nil {
		structTypes = append(structTypes, "client")
	}
	if entry.Binary != nil {
		structTypes = append(structTypes, "binary")
	}
	if entry.Asset != nil {
		structTypes = append(structTypes, "asset")
	}
	if entry.Geo != nil {
		structTypes = append(structTypes, "geo")
	}
	if entry.LDAP != nil {
		structTypes = append(structTypes, "ldap")
	}
	if entry.PDNS != nil {
		structTypes = append(structTypes, "pdns")
	}

	// If no structure types were found, add "unknown" as a fallback
	if len(structTypes) == 0 {
		structTypes = append(structTypes, "unknown")
	}

	return structTypes
}

// ------------------------------------------------Helper functions to parse nested data for each endpoint in FAKEula----------------------------------------------
func parseOil(entryMap map[string]interface{}) *OilInfo {
	oil := &OilInfo{
		Timestamp:     getString(entryMap, "timestamp"),
		UserPrincipal: getString(entryMap, "userPrincipalName"),
		DisplayName:   getString(entryMap, "displayName"),
		ClientIP:      getString(entryMap, "clientIp"),
		ClientASNOrg:  getString(entryMap, "clientAsOrg"),
		EventType:     getString(entryMap, "eventType"),
		Outcome:       getString(entryMap, "outcome"),
		Message:       getString(entryMap, "message"),
	}

	return oil
}

func parseClient(entryMap map[string]interface{}) *ClientInfo {
	// Check if the "client" field exists and is a map
	if client, ok := entryMap["client"].(map[string]interface{}); ok {
		// Return a new ClientInfo struct populated with data
		return &ClientInfo{
			AsOrg: getString(client, "as_org"),
			ASN:   getInt(client, "asn"),
			IP:    getString(client, "ip"),
		}
	}
	return nil
}

func parseBinary(entryMap map[string]interface{}) *BinaryInfo {
	if binaryData, ok := entryMap["binary"].(map[string]interface{}); ok {
		binary := &BinaryInfo{}

		// Try to extract file information
		if file, ok := binaryData["file"].(map[string]interface{}); ok {
			// Get filename
			binary.Filename = getString(file, "name")

			// Get accessed timestamp
			binary.Accessed = getString(file, "accessed")

			// Extract hash information
			if hash, ok := file["hash"].(map[string]interface{}); ok {
				binary.MD5 = getString(hash, "md5")
				// SHA256 might also be in the hash object if available, couldn't tell if it was or not in the FAKEula readme
				binary.SHA256 = getString(hash, "sha256")
			}

			// Extract host information
			if hosts, ok := file["hosts"].([]interface{}); ok {
				// Initialize a slice to store host names
				hostNames := make([]string, 0, len(hosts))
				// Iterate through each host
				for _, h := range hosts {
					if host, ok := h.(map[string]interface{}); ok {
						// Iterate through each host
						hostName := getString(host, "name")
						// Add it to the slice if it's not empty
						if hostName != "" {
							hostNames = append(hostNames, hostName)
						}
					}
				}
				binary.Hosts = hostNames
			}

			// Extract code signature information
			if signature, ok := file["code_signature"].(map[string]interface{}); ok {
				if exists, ok := signature["exists"].(bool); ok {
					binary.CodeSigned = exists
				}
			}
		}

		// Extract URL from labels if available
		if labels, ok := binaryData["labels"].(map[string]interface{}); ok {
			binary.URL = getString(labels, "url")
		}

		return binary
	}
	return nil
}

func parseAsset(entryMap map[string]interface{}) *AssetInfo {
	// First check if "asset" exists in the entry map
	if assetData, ok := entryMap["asset"].(map[string]interface{}); ok {
		asset := &AssetInfo{}

		// Try to extract host info
		if host, ok := assetData["host"].(map[string]interface{}); ok {
			asset.Name = getString(host, "name")
			asset.IP = getString(host, "ip")
		}

		// Try to extract platform info
		if platform, ok := assetData["platform"].(map[string]interface{}); ok {
			asset.PlatformName = getString(platform, "name")

			// Extract platform owner
			if owner, ok := platform["owner"].(map[string]interface{}); ok {
				asset.PlatformOwner = getString(owner, "full_name")
			}

			// Extract executive info
			if executive, ok := platform["executive"].(map[string]interface{}); ok {
				asset.Executive = getString(executive, "full_name")
			}
		}

		// Try to extract stack info
		if stack, ok := assetData["stack"].(map[string]interface{}); ok {
			asset.StackName = getString(stack, "name")

			// Extract stack owner
			if owner, ok := stack["owner"].(map[string]interface{}); ok {
				asset.StackOwner = getString(owner, "full_name")
			}
		}

		// Try to extract event timestamps
		if event, ok := assetData["event"].(map[string]interface{}); ok {
			asset.Created = getString(event, "created")
			asset.Updated = getString(event, "updated")
		}

		return asset
	}
	return nil
}

func parseGeo(entryMap map[string]interface{}) *GeoInfo {
	if geoData, ok := entryMap["geo"].(map[string]interface{}); ok {
		geo := &GeoInfo{}

		// Extract country and city info
		geo.CountryCode = getString(geoData, "country_iso_code")
		geo.CountryName = getString(geoData, "country_name")
		geo.City = getString(geoData, "city")
		geo.Latitude = getFloat(geoData, "latitude")
		geo.Longitude = getFloat(geoData, "longitude")

		// Try to extract AS info if available
		if as, ok := geoData["as"].(map[string]interface{}); ok {
			geo.ASNumber = getString(as, "number")

			// Extract organization name
			if org, ok := as["organization"].(map[string]interface{}); ok {
				geo.ASOrg = getString(org, "name")
			}
		}

		return geo
	}
	return nil
}

func parseLdap(entryMap map[string]interface{}) *LdapInfo {
	if ldap, ok := entryMap["ldap"].(map[string]interface{}); ok {
		return &LdapInfo{
			Email:       getString(ldap, "email"),
			FullName:    getString(ldap, "fullName"),
			Name:        getString(ldap, "name"),
			Title:       getString(ldap, "title"),
			CompanyName: getString(ldap, "companyName"),
			Phone:       getString(ldap, "phone"),
			Mobile:      getString(ldap, "mobile"),
			Created:     getString(ldap, "created"),
			Manager:     getString(ldap, "manager"),
			Age:         getString(ldap, "age"),
		}
	}
	return nil
}

func parsePDNS(entryMap map[string]interface{}) *PDNSInfo {
	if pdnsData, ok := entryMap["pdns"].(map[string]interface{}); ok {
		pdns := &PDNSInfo{
			Answers: []DNSAnswer{},
		}

		// Try to extract DNS answers
		if dns, ok := pdnsData["dns"].(map[string]interface{}); ok {
			// Iterate through each answer
			if answers, ok := dns["answers"].([]interface{}); ok {
				for _, a := range answers {
					if answer, ok := a.(map[string]interface{}); ok {
						// Create a new DNSAnswer struct and populate it
						dnsAnswer := DNSAnswer{
							Data:  getString(answer, "data"),
							Name:  getString(answer, "name"),
							Type:  getString(answer, "type"),
							Count: getInt(answer, "count"),
						}

						// Extract event times
						if event, ok := answer["event"].(map[string]interface{}); ok {
							dnsAnswer.Start = getString(event, "start")
							dnsAnswer.End = getString(event, "end")
						}

						// Add this answer to the slice
						pdns.Answers = append(pdns.Answers, dnsAnswer)
					}
				}
			}
		}

		return pdns
	}
	return nil
}

//-----------------------------------------------General Helper functions---------------------------------------------------------------------

// getSource determines the data source type for an entry by checking which specific data fields are present in the entry map
func getSource(entryMap map[string]interface{}) string {
	// 1. Check direct known keys
	for key := range entryMap {
		switch key {
		case "client", "binary", "asset", "geo", "ldap", "pdns":
			return key
		}
	}

	// 2. Check megaoil pipeline tag
	if megaoil, ok := entryMap["megaoil"].(map[string]interface{}); ok {
		if pipeline, ok := megaoil["pipeline"].(string); ok {
			return pipeline
		}
	}

	// 3. Check for fallback event.module
	if event, ok := entryMap["event"].(map[string]interface{}); ok {
		if mod, ok := event["module"].(string); ok {
			return mod
		}
	}

	// 4. Last resort: use "oil" field
	if oil, ok := entryMap["oil"].(string); ok && oil != "" {
		return oil
	}

	return "unknown"
}

// getString safely extracts a string value from a map using the provided key
// Returns empty string if the key doesn't exist or isn't a string
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

// getFloat safely extracts a float64 value from a map using the provided key
// Returns 0.0 if the key doesn't exist or isn't a float64
func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0.0
}

// getInt safely extracts an integer value from a map using the provided key
// I read In JSON, numbers are usually decoded as float64, so this converts to int
// Returns 0 if the key doesn't exist or isn't a number
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

// Print results cache for debugging
func PrintResultsCache() {
	fmt.Println("=== Printing resultsCache ===")
	for key, value := range resultsCache {
		fmt.Printf("Key: %s\nValue: %+v\n\n", key, value)
	}
	fmt.Println("=== End of resultsCache ===")
}
