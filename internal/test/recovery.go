package test

import (
    "github.com/wyw14/solo-037/internal/journal"
    "github.com/wyw14/solo-037/internal/model"
)

func Restore(records []journal.Record, run model.RunID) model.Operation {
    decision := journal.Recover(records, run)
    return model.Operation{ID: run, Stage: decision.Stage}
}

func CanRemoveModel(records []journal.Record, run model.RunID) bool {
    return journal.Recover(records, run).CanClear
}

func EvidenceCount(records []journal.Record, run model.RunID) int {
    return journal.Replay(records, run).Evidence
}
