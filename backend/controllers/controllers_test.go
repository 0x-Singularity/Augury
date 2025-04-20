package controllers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/0x-Singularity/Augury/controllers"
	"github.com/kylelemons/godebug/pretty"
)

// fakeFakeula spins up a local HTTP server that pretends to be the FAKEula API for testing.
// It answers every request with a minimal JSON
func fakeFakeula() *httptest.Server {
	mux := http.NewServeMux()

	// Generic catch‑all – return an empty data array
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []any{},
		})
	})

	// /extract needs to return an IOC so ExtractFromText can find it
	mux.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []any{
				map[string]any{
					"threat": map[string]any{
						"indicator": map[string]any{
							"description": "malicious.com",
						},
					},
				},
			},
		})
	})

	return httptest.NewServer(mux)
}

// helper returns a recorder/response JSON body and status.
func performRequest(h http.HandlerFunc, method, target string, body []byte) (*httptest.ResponseRecorder, map[string]any, error) {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	var decoded map[string]any
	if rr.Code == http.StatusOK {
		if err := json.NewDecoder(rr.Body).Decode(&decoded); err != nil {
			return rr, nil, err
		}
	}
	return rr, decoded, nil
}

func TestQueryPDNS_Success(t *testing.T) {
	server := fakeFakeula()
	defer server.Close()

	// point controller code at the fake server
	os.Setenv("FAKEULA_API_URL", server.URL+"/")
	os.Setenv("FAKEULA_USER", "user")
	os.Setenv("FAKEULA_PASS", "pass")
	os.Setenv("AUGURY_SKIP_DB", "1")

	rr, body, err := performRequest(controllers.QueryPDNS, http.MethodGet, "/pdns?ioc=example.com", nil)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if _, ok := body["data"]; !ok {
		t.Fatalf("expected 'data' field in response JSON")
	}
	t.Logf("TestQueryPDNA_Succesresponse body: %+v\n", body)
}

func TestQueryPDNS_MissingIOC(t *testing.T) {
	rr := httptest.NewRecorder()
	controllers.QueryPDNS(rr, httptest.NewRequest(http.MethodGet, "/pdns", nil))

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing ioc query‑param, got %d", rr.Code)
	}
	t.Logf("TestQueryPDNS_MissingIOC passed: %d", rr.Code)
}

func TestExtractFromText_Integration(t *testing.T) {
	server := fakeFakeula()
	defer server.Close()

	os.Setenv("FAKEULA_API_URL", server.URL+"/")
	os.Setenv("FAKEULA_USER", "user")
	os.Setenv("FAKEULA_PASS", "pass")
	os.Setenv("AUGURY_SKIP_DB", "1")

	text := []byte("visit http://malicious.com for more info")
	rr, body, err := performRequest(controllers.ExtractFromText, http.MethodPost, "/extract", text)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("response missing data map")
	}

	if _, present := data["malicious.com"]; !present {
		t.Fatalf("expected key for IOC 'malicious.com' in data map")
	}
	log.Println("TestExtractFromText_Integration passed")
	t.Log("response body:", pretty.Sprint(body))
}
