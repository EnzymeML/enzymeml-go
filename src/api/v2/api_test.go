package v2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
	"github.com/EnzymeML/enzymeml-go/src/database/v2"
)

func TestV2Routes_AllCollectionsRegistered(t *testing.T) {
	handler := newTestHandler(t)

	paths := []string{
		"/v2/documents",
		"/v2/creators",
		"/v2/vessels",
		"/v2/proteins",
		"/v2/complexes",
		"/v2/small-molecules",
		"/v2/reactions",
		"/v2/reaction-elements",
		"/v2/modifier-elements",
		"/v2/equations",
		"/v2/variables",
		"/v2/parameters",
		"/v2/measurements",
		"/v2/measurement-data",
		"/v2/unit-definitions",
		"/v2/base-units",
		"/v2/docs",
		"/v2/openapi.json",
	}

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s expected %d, got %d: %s", path, http.StatusOK, rec.Code, rec.Body.String())
		}
	}
}

func TestV2Routes_HumaDocsAndOpenAPI(t *testing.T) {
	handler := newTestHandler(t)

	docsResp := doRequest(t, handler, http.MethodGet, "/v2/docs", nil)
	if docsResp.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, docsResp.Code, docsResp.Body.String())
	}
	if !strings.Contains(strings.ToLower(docsResp.Body.String()), "openapi") {
		t.Fatalf("expected docs HTML to reference OpenAPI")
	}

	openapiResp := doRequest(t, handler, http.MethodGet, "/v2/openapi.json", nil)
	if openapiResp.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, openapiResp.Code, openapiResp.Body.String())
	}
	var spec map[string]interface{}
	if err := json.Unmarshal(openapiResp.Body.Bytes(), &spec); err != nil {
		t.Fatalf("failed to decode OpenAPI response: %v", err)
	}
	if spec["openapi"] == nil {
		t.Fatalf("expected openapi field in response")
	}
}

func TestV2Routes_CRUDCreatorAndProtein(t *testing.T) {
	handler := newTestHandler(t)

	creatorID := createResource(t, handler, "/v2/creators", map[string]interface{}{
		"given_name":  "Ada",
		"family_name": "Lovelace",
		"mail":        "ada@example.com",
	})
	if creatorID != float64(1) {
		t.Fatalf("expected creator id 1, got %#v", creatorID)
	}

	updateResource(t, handler, "/v2/creators/1", map[string]interface{}{
		"given_name":  "Augusta Ada",
		"family_name": "Lovelace",
		"mail":        "ada@example.com",
	})
	getRec := doRequest(t, handler, http.MethodGet, "/v2/creators/1", nil)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, getRec.Code, getRec.Body.String())
	}

	proteinID := createResource(t, handler, "/v2/proteins", map[string]interface{}{
		"id":       "p1",
		"name":     "Protein A",
		"constant": false,
	})
	if proteinID != "p1" {
		t.Fatalf("expected protein id p1, got %#v", proteinID)
	}

	updateResource(t, handler, "/v2/proteins/p1", map[string]interface{}{
		"id":       "ignored-on-update",
		"name":     "Protein A2",
		"constant": true,
	})
	getProteinRec := doRequest(t, handler, http.MethodGet, "/v2/proteins/p1", nil)
	if getProteinRec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, getProteinRec.Code, getProteinRec.Body.String())
	}

	deleteRec := doRequest(t, handler, http.MethodDelete, "/v2/proteins/p1", nil)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d: %s", http.StatusNoContent, deleteRec.Code, deleteRec.Body.String())
	}
	missingRec := doRequest(t, handler, http.MethodGet, "/v2/proteins/p1", nil)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d: %s", http.StatusNotFound, missingRec.Code, missingRec.Body.String())
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "v2_api_test.db")
	manager, err := database.NewDBManager(dbPath, v2Models())
	if err != nil {
		t.Fatalf("failed to create DB manager: %v", err)
	}
	t.Cleanup(func() { _ = manager.Close() })

	handler, err := NewHandler(manager)
	if err != nil {
		t.Fatalf("failed to create v2 handler: %v", err)
	}
	return handler
}

func v2Models() []interface{} {
	return []interface{}{
		&enzymeml_v2.EnzymeMLDocument{},
		&enzymeml_v2.Creator{},
		&enzymeml_v2.Vessel{},
		&enzymeml_v2.Protein{},
		&enzymeml_v2.Complex{},
		&enzymeml_v2.SmallMolecule{},
		&enzymeml_v2.Reaction{},
		&enzymeml_v2.ReactionElement{},
		&enzymeml_v2.ModifierElement{},
		&enzymeml_v2.Equation{},
		&enzymeml_v2.Variable{},
		&enzymeml_v2.Parameter{},
		&enzymeml_v2.Measurement{},
		&enzymeml_v2.MeasurementData{},
		&enzymeml_v2.UnitDefinition{},
		&enzymeml_v2.BaseUnit{},
	}
}

func createResource(t *testing.T, handler http.Handler, path string, payload interface{}) interface{} {
	t.Helper()

	rec := doRequest(t, handler, http.MethodPost, path, payload)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST %s expected %d, got %d: %s", path, http.StatusCreated, rec.Code, rec.Body.String())
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode POST %s response: %v", path, err)
	}
	return body["id"]
}

func updateResource(t *testing.T, handler http.Handler, path string, payload interface{}) {
	t.Helper()

	rec := doRequest(t, handler, http.MethodPut, path, payload)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT %s expected %d, got %d: %s", path, http.StatusOK, rec.Code, rec.Body.String())
	}
}

func doRequest(t *testing.T, handler http.Handler, method string, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to encode payload: %v", err)
		}
		body = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}
