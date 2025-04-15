package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	// Test Oil lookup
	//testOilLookup(t)

	// Test Client lookup
	//testClientLookup(t)

	// Test Process (base CBR) lookup
	testProcessLookup(t)

	// Test Host lookup
	//testHostLookup(t)

	// Test Binary lookup
	//testBinaryLookup(t)

	// Test Asset lookup
	//testAssetLookup(t)

	// Test Geo lookup
	//testGeoLookup(t)

	// Test LDAP lookup
	//testLdapLookup(t)

	// Test PDNS lookup
	//testPdnsLookup(t)
}

// Print function to send results of parsing to the console
func printResults(result map[string]map[string][]FakeulaEntry) {
	for source, sourceMap := range result {
		fmt.Printf("Source: %s\n", source)
		for structType, entries := range sourceMap {
			fmt.Printf("  Structure Type: %s\n", structType)
			for i, entry := range entries {
				fmt.Printf("    Entry %d: %+v\n", i+1, entry)

				// Print Oil data if present
				if entry.Oil != nil {
					fmt.Printf("      Oil Data: %+v\n", *entry.Oil)
				}

				// Print Client data if present
				if entry.Client != nil {
					fmt.Printf("      Client Data: %+v\n", *entry.Client)
				}

				// Print Process data if present
				if entry.Process != nil {
					fmt.Printf("      Process Data: %v\n", *&entry.Process)
				}

				// Print Host data if present
				if entry.Host != nil {
					fmt.Printf("      Host Data: %+v\n", *entry.Host)
				}

				// Print Binary data if present
				if entry.Binary != nil {
					fmt.Printf("      Binary Data: %+v\n", *entry.Binary)
				}

				// Print Asset data if present
				if entry.Asset != nil {
					fmt.Printf("      Asset Data: %+v\n", *entry.Asset)
				}

				// Print Geo data if present
				if entry.Geo != nil {
					fmt.Printf("      Geo Data: %+v\n", *entry.Geo)
				}

				// Print LDAP data if present
				if entry.LDAP != nil {
					fmt.Printf("      LDAP Data: %+v\n", *entry.LDAP)
				}

				// Print PDNS data if present
				if entry.PDNS != nil {
					fmt.Printf("      PDNS Data: %+v\n", *entry.PDNS)
					fmt.Printf("      PDNS Answers: %d records\n", len(entry.PDNS.Answers))
				}
			}
		}
	}
}

// Using Azure sample data for testing
func testOilLookup(t *testing.T) {
	oilJSON := `{
  "data": [
    {
      "callerIpAddress": "1.2.3.4",
      "coxAccountName": "abob",
      "userPrincipalName": "alice.bob@example.com",
      "userDisplayName": "Alice Bob",
      "displayName": "laptop1",
      "client": {
        "as_org": "ASN-ACME",
        "ip": "1.2.3.4",
        "asn": 1234
      },
      "timestamp": "2025-01-23T21:15:51.439Z",
      "key": "1.2.3.4",
      "oil": "azure"
    }
  ]
}`

	var oilResponse map[string]interface{}
	json.Unmarshal([]byte(oilJSON), &oilResponse)

	result := FormatFakeulaResponse(oilResponse)

	printResults(result.Data)
}

func testClientLookup(t *testing.T) {
	clientJSON := `{
  "data": [
    {
      "callerIpAddress": "1.2.3.4",
      "coxAccountName": "abob",
      "userPrincipalName": "alice.bob@example.com",
      "userDisplayName": "Alice Bob",
      "displayName": "laptop1",
      "client": {
        "as_org": "ASN-ACME",
        "ip": "1.2.3.4",
        "asn": 1234
      },
      "timestamp": "2025-01-23T21:15:51.439Z",
      "key": "1.2.3.4",
      "oil": "azure"
    }
  ]
}`

	var clientResponse map[string]interface{}
	json.Unmarshal([]byte(clientJSON), &clientResponse)

	result := FormatFakeulaResponse(clientResponse)

	printResults(result.Data)
}

func testProcessLookup(t *testing.T) {
	processJSON := `{
  "data": [
    {
      "process": {
        "command_line": "/usr/local/java/java_base/bin/java -Dp=executionserver -server -d64 -verbose:gc",
        "entity_id": "00003094-0000-13ad-01d8-6a476fb04ab2",
        "executable": "/bw/local/java/jdk1.8.0_312/bin/java",
        "name": "java",
        "pid": 5037,
        "start": "2022-05-17T23:40:05.647Z",
        "uptime": 53818,
        "hash": {
          "md5": "fb8b6d549055579989a7184077408342"
        },
        "parent": {
          "name": "bash",
          "pid": 5028,
          "entity_id": "00003094-0000-13a4-01d8-6a476faec8c6-000000000001"
        },
        "user": {
          "name": "alice",
          "id": null
        },
        "host": {
          "name": "host1",
          "type": "workstation",
          "ip": [
            "192.168.0.1"
          ],
          "os": {
            "family": "linux"
          }
        },
        "code_signature": {
          "exists": false
        }
      },
      "labels": {
        "url": "https://cbr.example.com/#/analyze/00003094-0000-13ad-01d8-6a476fb04ab2/1652830876949?cb.legacy_5x_mode=false"
      }
    }
  ]
}`

	var processResponse map[string]interface{}
	json.Unmarshal([]byte(processJSON), &processResponse)

	result := FormatFakeulaResponse(processResponse)

	// Print and verify the result
	// Similar to the client lookup verification
	printResults(result.Data)
}

func testHostLookup(t *testing.T) {
	// sample host response
	hostJSON := `{
  "data": [
    {
      "sensor": {
        "hostname": "host1",
        "id": 78065,
        "ip": [
          "192.168.0.1"
        ],
        "mac": [
          "00:00:00:00:00:01"
        ],
        "name": "host1.example.com",
        "uptime": 593924,
        "os": {
          "full": "Windows 10 Enterprise, 64-bit",
          "version": "007.003.000.18311"
        }
      },
      "labels": {
        "url": "https://cbr.example.com/#/host/78065"
      }
    }
  ]
}`

	var hostResponse map[string]interface{}
	json.Unmarshal([]byte(hostJSON), &hostResponse)

	result := FormatFakeulaResponse(hostResponse)

	// Print and verify the result
	printResults(result.Data)
}

func testBinaryLookup(t *testing.T) {
	// Sample binary response
	binaryJSON := `{
  "data": [
    {
      "file": {
        "hash": {
          "md5": "F88ADB10AB5313D4FA33416F6F5FB4FF"
        },
        "name": "ysoserial.exe",
        "accessed": "2022-04-27T11:50:32.029Z",
        "hosts": [
          {
            "name": "host1",
            "id": "17864"
          }
        ],
        "code_signature": {
          "exists": false
        }
      },
      "labels": {
        "url": "https://cbr.example.com/#/binary/F88ADB10AB5313D4FA33416F6F5FB4FF"
      }
    }
  ]
}`

	var binaryResponse map[string]interface{}
	json.Unmarshal([]byte(binaryJSON), &binaryResponse)

	result := FormatFakeulaResponse(binaryResponse)

	// Print and verify the result
	printResults(result.Data)
}

func testAssetLookup(t *testing.T) {
	// Sample asset response
	assetJSON := `{
  "data": [
    {
      "host": {
        "name": "SERVER.EXAMPLE.COM",
        "ip": "10.0.0.1"
      },
      "platform": {
        "name": "Security Investigator",
        "owner": {
          "full_name": "Alice"
        },
        "executive": {
          "full_name": "Bob"
        }
      },
      "stack": {
        "name": "Cyber Defense",
        "owner": {
          "full_name": "Charlie"
        }
      },
      "event": {
        "created": "2024-05-07T00:00:00Z",
        "updated": "2025-01-15T00:00:00Z"
      }
    }
  ]
}`

	var assetResponse map[string]interface{}
	json.Unmarshal([]byte(assetJSON), &assetResponse)

	result := FormatFakeulaResponse(assetResponse)
	printResults(result.Data)
}

func testGeoLookup(t *testing.T) {
	// Sample geo response
	geoJSON := `{
  "data": [
    {
      "host": {
        "ip": [
          "1.2.3.4"
        ]
      },
      "as": {
        "number": "",
        "organization": {
          "name": ""
        }
      },
      "geo": {
        "country_iso_code": "AU",
        "country_name": "Australia"
      }
    }
  ]
}`

	var geoResponse map[string]interface{}
	json.Unmarshal([]byte(geoJSON), &geoResponse)

	result := FormatFakeulaResponse(geoResponse)
	printResults(result.Data)
}

func testLdapLookup(t *testing.T) {
	// Sample LDAP response
	ldapJSON := `{
  "data": [
    {
      "user": {
        "email": "alice.bob@example.com",
        "full_name": "Alice Bob",
        "name": "abob",
        "title": "CISO",
        "company": "Example Corp",
        "phone": "+1 (555) 555-5555",
        "mobile": "+1 (555) 777-7777",
        "created": "2015-01-01 20:00:00+00:00",
        "manager": "CN=Bob Charlie (Example-Atlanta) bcharlie,OU=Users,OU=Standard Users,OU=Users and Computers,OU=Atlanta,OU=Example,DC=DOMAIN,DC=EXAMPLE,DC=com",
        "age": 8692
      }
    }
  ]
}`

	var ldapResponse map[string]interface{}
	json.Unmarshal([]byte(ldapJSON), &ldapResponse)

	result := FormatFakeulaResponse(ldapResponse)
	printResults(result.Data)
}

func testPdnsLookup(t *testing.T) {
	// Sample PDNS response
	pdnsJSON := `{
  "data": [
    {
      "dns": {
        "answers": [
          {
            "data": "1.2.3.4",
            "name": "a.internal-test-ignore.biz",
            "type": "A",
            "count": 1346,
            "event": {
              "start": "2019-11-06T22:54:18Z",
              "end": "2025-01-23T00:23:21Z"
            }
          },
          {
            "data": "1.2.3.4",
            "name": "b.internal-test-ignore.biz",
            "type": "A",
            "count": 1346,
            "event": {
              "start": "2019-11-06T22:54:18Z",
              "end": "2025-01-23T00:23:21Z"
            }
          },
          {
            "data": "1.2.3.4",
            "name": "ns1.37cw.com",
            "type": "A",
            "count": 1,
            "event": {
              "start": "2023-03-11T22:50:20Z",
              "end": "2023-03-11T22:50:20Z"
            }
          },
          {
            "data": "1.2.3.4",
            "name": "ns2.37cw.com",
            "type": "A",
            "count": 1,
            "event": {
              "start": "2023-03-11T22:50:20Z",
              "end": "2023-03-11T22:50:20Z"
            }
          },
          {
            "data": "1.2.3.4",
            "name": "ns1.45mov.com",
            "type": "A",
            "count": 1,
            "event": {
              "start": "2023-03-11T22:50:20Z",
              "end": "2023-03-11T22:50:20Z"
            }
          }
        ]
      }
    }
  ]
}`

	var pdnsResponse map[string]interface{}
	json.Unmarshal([]byte(pdnsJSON), &pdnsResponse)

	result := FormatFakeulaResponse(pdnsResponse)
	printResults(result.Data)
}
