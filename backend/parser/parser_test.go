package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	// Test Client lookup
	testClientLookup(t)

	// Test Binary lookup
	testBinaryLookup(t)

	// Test Asset lookup
	testAssetLookup(t)

	// Test Geo lookup
	testGeoLookup(t)

	// Test LDAP lookup
	testLdapLookup(t)

	// Test PDNS lookup
	testPdnsLookup(t)
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

// Print function to send results of parsing to the console
func printResults(result map[string]map[string][]FakeulaEntry) {
	for oil, sourceMap := range result {
		fmt.Printf("Oil: %s\n", oil)
		for source, entries := range sourceMap {
			fmt.Printf("  Source: %s\n", source)
			for _, entry := range entries {
				fmt.Printf("    Entry: %+v\n", entry)
			}
		}
	}
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
	// Similar to the client lookup verification
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
