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

func TestAbortUnknownRunReportsError(t *testing.T) {
	runtime, err := service.OpenRuntime(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	server := NewRouter(runtime).Handler()
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/abort?run=run-does-not-exist", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("abort unknown run: %d %s", response.Code, response.Body)
	}
	if !strings.Contains(response.Body.String(), "run not found") {
		t.Fatalf("abort unknown run should explain the failure, got %s", response.Body)
	}
	incidents := httptest.NewRecorder()
	server.ServeHTTP(incidents, httptest.NewRequest(http.MethodGet, "/api/incidents", nil))
	if strings.Contains(incidents.Body.String(), "run-does-not-exist") {
		t.Fatalf("unknown run must not record an incident, got %s", incidents.Body)
	}
}

func TestAbortKnownRunKeepsIncidentAndStage(t *testing.T) {
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
	if operation.ID == "" {
		t.Fatal("missing ID")
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
	found := false
	for _, candidate := range listed.Data {
		if candidate.ID == operation.ID {
			found = true
			if candidate.Stage != model.StageAborted {
				t.Fatalf("stage = %s, want %s", candidate.Stage, model.StageAborted)
			}
			if candidate.UpdatedAt.Before(operation.CreatedAt) {
				t.Fatalf("updated_at %s not refreshed from %s", candidate.UpdatedAt, operation.CreatedAt)
			}
		}
	}
	if !found {
		t.Fatal("aborted operation missing from operations list")
	}
	incidents := httptest.NewRecorder()
	server.ServeHTTP(incidents, httptest.NewRequest(http.MethodGet, "/api/incidents", nil))
	var recorded struct {
		Data []model.Incident `json:"data"`
	}
	if err := json.Unmarshal(incidents.Body.Bytes(), &recorded); err != nil {
		t.Fatal(err)
	}
	for _, incident := range recorded.Data {
		if incident.Run == operation.ID {
			return
		}
	}
	t.Fatalf("abort incident missing, got %s", incidents.Body)
}
