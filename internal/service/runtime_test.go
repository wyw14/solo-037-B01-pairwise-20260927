package service

import (
	"errors"
	"testing"
	"time"

	"github.com/wyw14/solo-037/internal/model"
)

func TestAbortUnknownRunReportsError(t *testing.T) {
	runtime, err := OpenRuntime(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.Abort("missing-run", time.Now().UTC()); !errors.Is(err, ErrRunNotFound) {
		t.Fatalf("unknown run: got %v, want ErrRunNotFound", err)
	}
	if incidents := runtime.Incidents().([]model.Incident); len(incidents) != 0 {
		t.Fatalf("unknown run recorded an incident: %+v", incidents)
	}
}

func TestAbortKnownRunKeepsIncidentAndUpdatedStage(t *testing.T) {
	runtime, err := OpenRuntime(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	operation, err := runtime.Create(model.NewOperation("tow", "hull-x", "section-a", "carriage-1", time.Now().UTC()))
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.Abort(operation.ID, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	var aborted *model.Operation
	for _, candidate := range runtime.Operations().([]model.Operation) {
		if candidate.ID == operation.ID {
			found := candidate
			aborted = &found
		}
	}
	if aborted == nil || aborted.Stage != model.StageAborted {
		t.Fatalf("operation stage not updated: %+v", aborted)
	}
	incidents := runtime.Incidents().([]model.Incident)
	if len(incidents) != 1 || incidents[0].Run != operation.ID {
		t.Fatalf("incident record missing: %+v", incidents)
	}
}
