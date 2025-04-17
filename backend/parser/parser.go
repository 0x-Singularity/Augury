package parser

import (
	"encoding/json"
	"fmt"
	"log"
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
	Oil     *OilInfo     `json:"oil"`
	Client  *ClientInfo  `json:"client,omitempty"`
	Process *ProcessInfo `json:"process,omitempty"`
	Host    *HostInfo    `json:"host,omitempty"`
	Binary  *BinaryInfo  `json:"binary,omitempty"`
	Asset   *AssetInfo   `json:"asset,omitempty"`
	Geo     *GeoInfo     `json:"geo,omitempty"`
	LDAP    *LdapInfo    `json:"ldap,omitempty"`
	PDNS    *PDNSInfo    `json:"pdns,omitempty"`
}

// OilInfo
type OilInfo struct {
	Timestamp     string `json:"timestamp,omitempty"`
	UserPrincipal string `json:"userPrincipalName,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	ClientIP      string `json:"clientIp,omitempty"`
	ClientASNOrg  string `json:"clientAsOrg,omitempty"`
	EventType     string `json:"eventType,omitempty"`
	Outcome       string `json:"outcome,omitempty"`
	Message       string `json:"message,omitempty"`

	// Optional fields from Okta, Helios, Netflow, and Prisma

	// Okta
	DisplayMessage string `json:"displayMessage"`

	// Helios
	ObserverHostname                string   `json:"observerHostname,omitempty"`
	SuricataSignature               string   `json:"suricataSignature,omitempty"`
	SourcePort                      string   `json:"sourcePort,omitempty"`
	SourceThreatClassification      string   `json:"sourceThreatClassification,omitempty"`
	SourceThreatService             string   `json:"sourceThreatService,omitempty"`
	DestinationThreatClassification string   `json:"destinationThreatClassification,omitempty"`
	DestinationThreatService        string   `json:"destinationThreatService,omitempty"`
	Pipeline                        string   `json:"pipeline,omitempty"`
	Tags                            []string `json:"tags,omitempty"`

	// Prisma
	NetworkProtocol    string `json:"networkProtocol,omitempty"`
	Application        string `json:"application,omitempty"`
	RuleName           string `json:"ruleName,omitempty"`
	EventSequence      string `json:"eventSequence,omitempty"`
	EventAction        string `json:"eventAction,omitempty"`
	SourcePackets      string `json:"sourcePackets,omitempty"`
	SourceBytes        string `json:"sourceBytes,omitempty"`
	DestinationPackets string `json:"destinationPackets,omitempty"`
	DestinationBytes   string `json:"destinationBytes,omitempty"`
	Transport          string `json:"transport,omitempty"`

	// Netflow
	EventStart string `json:"eventStart,omitempty"`
	EventEnd   string `json:"eventEnd,omitempty"`

	// Common geo-as org info
	SourceASNOrg  string `json:"sourceASNOrg,omitempty"`
	SourceASN     string `json:"sourceASN,omitempty"`
	SourceCountry string `json:"sourceCountry,omitempty"`
	SourceCity    string `json:"sourceCity,omitempty"`

	// Destination info (can also apply to Helios/Prisma/Netflow)
	DestinationIP   string `json:"destinationIP,omitempty"`
	DestinationPort string `json:"destinationPort,omitempty"`
	DestinationASN  string `json:"destinationASN,omitempty"`
	DestinationOrg  string `json:"destinationOrg,omitempty"`
}

// ClientInfo represents network client information
type ClientInfo struct {
	AsOrg string `json:"as_org"`
	ASN   int    `json:"asn"`
	IP    string `json:"ip"`
}

// ProcessInfo struct to match base CBR response, which is nested under the key word "process"
type ProcessInfo struct {
	Name           string   `json:"name"`
	CommandLine    string   `json:"command_line"`
	EntityID       string   `json:"entity_id"`
	Executable     string   `json:"executable"`
	PID            int      `json:"pid"`
	Start          string   `json:"start"`
	Uptime         int      `json:"uptime"`
	ParentName     string   `json:"parent_name"`
	ParentPID      int      `json:"parent_pid"`
	ParentEntityID string   `json:"parent_entity_id"`
	UserName       string   `json:"user_name"`
	HostName       string   `json:"host_name"`
	HostType       string   `json:"host_type"`
	HostIPs        []string `json:"host_ips"`
	HostOS         string   `json:"host_os"`
	CodeSigned     bool     `json:"code_signed"`
	URL            string   `json:"url"`
}

// HostInfo struct to match CBR Host response
type HostInfo struct {
	Hostname string   `json:"hostname"`
	Name     string   `json:"name"`
	ID       int      `json:"id"`
	IPs      []string `json:"ips"`
	MACs     []string `json:"macs"`
	Uptime   int      `json:"uptime"`
	OSFull   string   `json:"os_full"`
	OSVer    string   `json:"os_version"`
	URL      string   `json:"url"`
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
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	//City        string `json:"city"`
	//Latitude    float64 `json:"latitude"`
	//Longitude   float64 `json:"longitude"`
	ASNumber string `json:"asNumber"`
	ASOrg    string `json:"asOrg"`
	IP       string `json:"ip"`
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
					Oil:     parseOil(entryMap),
					Client:  parseClient(entryMap),
					Process: parseProcess(entryMap),
					Host:    parseHost(entryMap),
					Binary:  parseBinary(entryMap),
					Asset:   parseAsset(entryMap),
					Geo:     parseGeo(entryMap),
					LDAP:    parseLdap(entryMap),
					PDNS:    parsePDNS(entryMap),
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
					parsedData[source][structType] = append(parsedData[source][structType], parsedEntry)

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
	if entry.Process != nil {
		structTypes = append(structTypes, "process")
	}
	if entry.Host != nil {
		structTypes = append(structTypes, "host")
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
		// Added @timestamp fallback for Helios lookup
		Timestamp: func() string {
			if ts := getString(entryMap, "timestamp"); ts != "" {
				return ts
			}
			return getString(entryMap, "@timestamp")
		}(),
		UserPrincipal: getString(entryMap, "userPrincipalName"),
		DisplayName:   getString(entryMap, "displayName"),
		ClientIP:      getString(entryMap, "callerIpAddress"),
		ClientASNOrg:  "", // fallback set below
		EventType:     "",
		Outcome:       "",
		Message:       "",
	}

	// Client IP & ASN org
	if client, ok := entryMap["client"].(map[string]interface{}); ok {
		oil.ClientASNOrg = getString(client, "as_org")
		if ip := getString(client, "ipAddress"); ip != "" {
			oil.ClientIP = ip
		}
	}

	// actor block (Okta)
	if actor, ok := entryMap["actor"].(map[string]interface{}); ok {
		oil.UserPrincipal = getString(actor, "alternateId")
		oil.DisplayName = getString(actor, "displayName")
	}

	// displayMessage (Okta fallback)
	if msg := getString(entryMap, "displayMessage"); msg != "" {
		oil.Message = msg
		oil.DisplayMessage = msg
	}

	// event block
	if event, ok := entryMap["event"].(map[string]interface{}); ok {
		oil.EventType = getString(event, "type")
		oil.Outcome = getString(event, "outcome")
		if oil.Message == "" {
			oil.Message = getString(event, "message")
		}
		oil.EventSequence = getString(event, "sequence")
		oil.EventAction = getString(event, "action")
		oil.EventStart = getString(event, "start")
		oil.EventEnd = getString(event, "end")
	}

	// Suricata
	if suricata, ok := entryMap["Suricata"].(map[string]interface{}); ok {
		oil.SuricataSignature = getString(suricata, "Signature")
	}

	// observer.hostname
	if observer, ok := entryMap["observer"].(map[string]interface{}); ok {
		oil.ObserverHostname = getString(observer, "hostname")
	}

	// rule (Prisma)
	if rule, ok := entryMap["rule"].(map[string]interface{}); ok {
		oil.RuleName = getString(rule, "name")
	}

	// transport (Prisma)
	if network, ok := entryMap["network"].(map[string]interface{}); ok {
		oil.NetworkProtocol = getString(network, "protocol")
		oil.Application = getString(network, "application")
		oil.Transport = getString(network, "transport")
	}

	// source (Prisma)
	if source, ok := entryMap["source"].(map[string]interface{}); ok {
		oil.SourcePackets = getString(source, "packets")
		oil.SourceBytes = getString(source, "bytes")
	}

	// destination (Prisma)
	if dest, ok := entryMap["destination"].(map[string]interface{}); ok {
		oil.DestinationPackets = getString(dest, "packets")
		oil.DestinationBytes = getString(dest, "bytes")
	}

	// network info
	if network, ok := entryMap["network"].(map[string]interface{}); ok {
		oil.NetworkProtocol = getString(network, "protocol")
		oil.Application = getString(network, "application")
	}

	// source block
	if source, ok := entryMap["source"].(map[string]interface{}); ok {
		if ip := getString(source, "ip"); ip != "" && oil.ClientIP == "" {
			oil.ClientIP = ip
		}
		oil.SourcePort = getString(source, "port")

		if geo, ok := source["geo"].(map[string]interface{}); ok {
			oil.SourceCountry = getString(geo, "country_iso_code")
			oil.SourceCity = getString(geo, "city_name")
		}
		if as, ok := source["as"].(map[string]interface{}); ok {
			if number, ok := as["number"]; ok && number != nil {
				oil.SourceASN = fmt.Sprintf("%v", number)
			}
			if org, ok := as["organization"].(map[string]interface{}); ok {
				oil.SourceASNOrg = getString(org, "name")
			}
		}
		if threat, ok := source["threat"].(map[string]interface{}); ok {
			if indicator, ok := threat["indicator"].(map[string]interface{}); ok {
				oil.SourceThreatClassification = getString(indicator, "Classification")
				oil.SourceThreatService = getString(indicator, "Service_Name")
			}
		}
	}

	// destination block
	if destination, ok := entryMap["destination"].(map[string]interface{}); ok {
		oil.DestinationIP = getString(destination, "ip")
		oil.DestinationPort = getString(destination, "port")

		if as, ok := destination["as"].(map[string]interface{}); ok {
			if number, ok := as["number"]; ok && number != nil {
				oil.DestinationASN = fmt.Sprintf("%v", number)
			}
			if org, ok := as["organization"].(map[string]interface{}); ok {
				oil.DestinationOrg = getString(org, "name")
			}
		}
		if threat, ok := destination["threat"].(map[string]interface{}); ok {
			if indicator, ok := threat["indicator"].(map[string]interface{}); ok {
				oil.DestinationThreatClassification = getString(indicator, "Classification")
				oil.DestinationThreatService = getString(indicator, "Service_Name")
			}
		}
	}

	// megaoil pipeline
	if megaoil, ok := entryMap["megaoil"].(map[string]interface{}); ok {
		oil.Pipeline = getString(megaoil, "pipeline")
	}

	// tags (string slice)
	if tagsRaw, ok := entryMap["tags"].([]interface{}); ok {
		for _, tag := range tagsRaw {
			if s, ok := tag.(string); ok {
				oil.Tags = append(oil.Tags, s)
			}
		}
	}

	// minimal data check
	if oil.Timestamp == "" &&
		oil.ClientIP == "" &&
		oil.UserPrincipal == "" &&
		oil.DisplayName == "" &&
		oil.ObserverHostname == "" &&
		oil.Message == "" &&
		oil.SuricataSignature == "" &&
		oil.EventStart == "" &&
		oil.NetworkProtocol == "" &&
		oil.SourceASNOrg == "" {
		return nil
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

// CBR has three different types of responses, Host, Process, and Binary

// Process (base CBR response)
func parseProcess(entryMap map[string]interface{}) *ProcessInfo {
	log.Println("Parsing process info")
	log.Println(entryMap)
	if processMap, ok := entryMap["process"].(map[string]interface{}); ok {
		process := &ProcessInfo{
			CommandLine: getString(processMap, "command_line"),
			EntityID:    getString(processMap, "entity_id"),
			Executable:  getString(processMap, "executable"),
			Name:        getString(processMap, "name"),
			PID:         getInt(processMap, "pid"),
			Start:       getString(processMap, "start"),
			Uptime:      getInt(processMap, "uptime"),
		}

		// Parent info
		if parent, ok := processMap["parent"].(map[string]interface{}); ok {
			process.ParentName = getString(parent, "name")
			process.ParentPID = getInt(parent, "pid")
			process.ParentEntityID = getString(parent, "entity_id")
		}

		// User info
		if user, ok := processMap["user"].(map[string]interface{}); ok {
			process.UserName = getString(user, "name")
		}

		// Host info
		if host, ok := processMap["host"].(map[string]interface{}); ok {
			process.HostName = getString(host, "name")
			process.HostType = getString(host, "type")

			// Check if there are multiple ips, then append them
			if ips, ok := host["ip"].([]interface{}); ok {
				for _, ip := range ips {
					if ipStr, ok := ip.(string); ok {
						process.HostIPs = append(process.HostIPs, ipStr)
					}
				}
			}

			if os, ok := host["os"].(map[string]interface{}); ok {
				process.HostOS = getString(os, "family")
			}
		}

		// Code signature
		if signature, ok := processMap["code_signature"].(map[string]interface{}); ok {
			if exists, ok := signature["exists"].(bool); ok {
				process.CodeSigned = exists
			}
		}

		// URL from labels
		if labels, ok := entryMap["labels"].(map[string]interface{}); ok {
			process.URL = getString(labels, "url")
		}

		if process.Name != "" || process.Executable != "" {
			return process
		}
	}
	return nil
}

// Host
func parseHost(entryMap map[string]interface{}) *HostInfo {
	// "sensor" is the keyword in the JSON response that Host info is nested under
	if sensor, ok := entryMap["sensor"].(map[string]interface{}); ok {
		host := &HostInfo{
			Hostname: getString(sensor, "hostname"),
			Name:     getString(sensor, "name"),
			ID:       getInt(sensor, "id"),
			Uptime:   getInt(sensor, "uptime"),
		}

		// IP addresses
		if ips, ok := sensor["ip"].([]interface{}); ok {
			for _, ip := range ips {
				if s, ok := ip.(string); ok {
					host.IPs = append(host.IPs, s)
				}
			}
		}

		// MAC addresses
		if macs, ok := sensor["mac"].([]interface{}); ok {
			for _, mac := range macs {
				if s, ok := mac.(string); ok {
					host.MACs = append(host.MACs, s)
				}
			}
		}

		// OS info
		if os, ok := sensor["os"].(map[string]interface{}); ok {
			host.OSFull = getString(os, "full")
			host.OSVer = getString(os, "version")
		}

		// Labels (URL)
		if labels, ok := entryMap["labels"].(map[string]interface{}); ok {
			host.URL = getString(labels, "url")
		}

		if host.Hostname != "" || host.Name != "" {
			return host
		}
	}
	return nil
}

// Binary
func parseBinary(entryMap map[string]interface{}) *BinaryInfo {
	// Try to extract file information
	// Skip looking for "binary" keyword, binary responses are nested under "file"
	if file, ok := entryMap["file"].(map[string]interface{}); ok {
		binary := &BinaryInfo{}

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
			for _, h := range hosts {
				if host, ok := h.(map[string]interface{}); ok {
					if name := getString(host, "name"); name != "" {
						binary.Hosts = append(binary.Hosts, name)
					}
				}
			}
		}

		// Extract code signature information
		if signature, ok := file["code_signature"].(map[string]interface{}); ok {
			if exists, ok := signature["exists"].(bool); ok {
				binary.CodeSigned = exists
			}
		}

		// Extract URL from labels if available
		if labels, ok := entryMap["labels"].(map[string]interface{}); ok {
			binary.URL = getString(labels, "url")
		}

		// Check if it's truly populated
		if binary.MD5 != "" || binary.Filename != "" {
			return binary
		}
	}
	return nil
}

func parseAsset(entryMap map[string]interface{}) *AssetInfo {
	asset := &AssetInfo{}

	// Try to extract host info
	if host, ok := entryMap["host"].(map[string]interface{}); ok {
		asset.Name = getString(host, "name")
		asset.IP = getString(host, "ip")
	}

	// Try to extract platform info
	if platform, ok := entryMap["platform"].(map[string]interface{}); ok {
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
	if stack, ok := entryMap["stack"].(map[string]interface{}); ok {
		asset.StackName = getString(stack, "name")

		// Extract stack owner
		if owner, ok := stack["owner"].(map[string]interface{}); ok {
			asset.StackOwner = getString(owner, "full_name")
		}
	}

	// Try to extract event timestamps
	if event, ok := entryMap["event"].(map[string]interface{}); ok {
		asset.Created = getString(event, "created")
		asset.Updated = getString(event, "updated")
	}

	// Return only if meaningful
	if asset.Name != "" || asset.IP != "" || asset.PlatformName != "" {
		return asset
	}

	return nil
}

func parseGeo(entryMap map[string]interface{}) *GeoInfo {
	if geoData, ok := entryMap["geo"].(map[string]interface{}); ok {
		geo := &GeoInfo{
			CountryCode: getString(geoData, "country_iso_code"),
			CountryName: getString(geoData, "country_name"),
			//City:        getString(geoData, "city"),
			//Latitude:    getFloat(geoData, "latitude"),
			//Longitude:   getFloat(geoData, "longitude"),
		}

		// Pull IP from "host.ip" if available
		if host, ok := entryMap["host"].(map[string]interface{}); ok {
			if ipList, ok := host["ip"].([]interface{}); ok && len(ipList) > 0 {
				if ipStr, ok := ipList[0].(string); ok {
					geo.IP = ipStr
				}
			}
		}

		// Pull ASN info from top-level "as" field
		if asData, ok := entryMap["as"].(map[string]interface{}); ok {
			geo.ASNumber = getString(asData, "number")

			if org, ok := asData["organization"].(map[string]interface{}); ok {
				geo.ASOrg = getString(org, "name")
			}
		}

		// Return only if country or IP exists
		if geo.CountryCode != "" || geo.IP != "" {
			return geo
		}
	}
	return nil
}

func parseLdap(entryMap map[string]interface{}) *LdapInfo {
	if user, ok := entryMap["user"].(map[string]interface{}); ok {
		ldap := &LdapInfo{
			Email:       getString(user, "email"),
			FullName:    getString(user, "full_name"),
			Name:        getString(user, "name"),
			Title:       getString(user, "title"),
			CompanyName: getString(user, "company"),
			Phone:       getString(user, "phone"),
			Mobile:      getString(user, "mobile"),
			Created:     getString(user, "created"),
			Manager:     getString(user, "manager"),
			Age:         fmt.Sprintf("%v", user["age"]),
		}

		// Basic check to avoid empty structs
		if ldap.Email != "" || ldap.Name != "" {
			return ldap
		}
	}
	return nil
}

func parsePDNS(entryMap map[string]interface{}) *PDNSInfo {
	// Try to extract DNS answers
	if dns, ok := entryMap["dns"].(map[string]interface{}); ok {
		pdns := &PDNSInfo{
			Answers: []DNSAnswer{},
		}

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
		if len(pdns.Answers) > 0 {
			return pdns
		}
	}

	return nil
}

//-----------------------------------------------General Helper functions---------------------------------------------------------------------

// getSource determines the data source type for an entry by checking which specific data fields are present in the entry map
func getSource(entryMap map[string]interface{}) string {
	// Attempt to determine the source based on the presence of known substructures.

	// Check client info
	if _, ok := entryMap["client"].(map[string]interface{}); ok {
		return "client"
	}

	//Check CBR responses

	//Check process info
	if _, ok := entryMap["process"].(map[string]interface{}); ok {
		return "process"
	}

	//Check host info
	if _, ok := entryMap["sensor"].(map[string]interface{}); ok {
		return "host"
	}

	// Check binary info (direct or nested under "file")
	if _, ok := entryMap["binary"].(map[string]interface{}); ok {
		return "binary"
	}
	if _, ok := entryMap["file"].(map[string]interface{}); ok {
		return "binary"
	}

	// Check asset info
	if _, ok := entryMap["asset"].(map[string]interface{}); ok {
		return "asset"
	}
	if _, ok := entryMap["host"].(map[string]interface{}); ok {
		if _, ok := entryMap["platform"].(map[string]interface{}); ok {
			return "asset"
		}
	}

	// Check geo info
	if _, ok := entryMap["geo"].(map[string]interface{}); ok {
		return "geo"
	}
	if _, ok := entryMap["as"].(map[string]interface{}); ok {
		if _, ok := entryMap["geo"].(map[string]interface{}); ok {
			return "geo"
		}
	}

	// Check LDAP info
	if _, ok := entryMap["ldap"].(map[string]interface{}); ok {
		return "ldap"
	}
	if _, ok := entryMap["user"].(map[string]interface{}); ok {
		return "ldap"
	}

	// Check PDNS info
	if _, ok := entryMap["pdns"].(map[string]interface{}); ok {
		return "pdns"
	}
	if dns, ok := entryMap["dns"].(map[string]interface{}); ok {
		if _, ok := dns["answers"].([]interface{}); ok {
			return "pdns"
		}
	}

	// Fallback to "oil" field if it exists and has a recognizable value
	if oil, ok := entryMap["oil"].(string); ok && oil != "" {
		return oil
	}

	// Try megaoil pipeline
	if megaoil, ok := entryMap["megaoil"].(map[string]interface{}); ok {
		if pipeline, ok := megaoil["pipeline"].(string); ok {
			return pipeline
		}
	}

	// Try event.module as a fallback
	if event, ok := entryMap["event"].(map[string]interface{}); ok {
		if mod, ok := event["module"].(string); ok {
			return mod
		}
	}

	// Default fallback
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
