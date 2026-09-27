package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wyw14/solo-037/internal/model"
	"github.com/wyw14/solo-037/internal/service"
)

func TestAbortUnknownRunReturnsNotFound(t *testing.T) {
	runtime, err := service.OpenRuntime(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	server := NewRouter(runtime).Handler()
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/abort?run=missing-run", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("unknown run: got %d %s, want 404", response.Code, response.Body)
	}
}

func TestAbortKnownRunKeepsIncidentAndUpdatedStage(t *testing.T) {
	runtime, err := service.OpenRuntime(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	server := NewRouter(runtime).Handler()
	created := httptest.NewRecorder()
	server.ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/api/operations", nil))
	if created.Code != http.StatusOK {
		t.Fatalf("create: %d %s", created.Code, created.Body)
	}
	var operation model.Operation
	if err := json.Unmarshal(created.Body.Bytes(), &operation); err != nil {
		t.Fatal(err)
	}
	aborted := httptest.NewRecorder()
	server.ServeHTTP(aborted, httptest.NewRequest(http.MethodPost, "/api/abort?run="+string(operation.ID), nil))
	if aborted.Code != http.StatusNoContent {
		t.Fatalf("abort: %d %s", aborted.Code, aborted.Body)
	}
	operations := httptest.NewRecorder()
	server.ServeHTTP(operations, httptest.NewRequest(http.MethodGet, "/api/operations", nil))
	var listed struct {
		Data []model.Operation `json:"data"`
	}
	if err := json.Unmarshal(operations.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	stageUpdated := false
	for _, candidate := range listed.Data {
		if candidate.ID == operation.ID {
			stageUpdated = candidate.Stage == model.StageAborted
		}
	}
	if !stageUpdated {
		t.Fatalf("operation stage not updated: %s", operations.Body)
	}
	incidents := httptest.NewRecorder()
	server.ServeHTTP(incidents, httptest.NewRequest(http.MethodGet, "/api/incidents", nil))
	if !strings.Contains(incidents.Body.String(), string(operation.ID)) {
		t.Fatalf("incident record missing: %s", incidents.Body)
	}
}
