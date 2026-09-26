package service

import (
    "github.com/wyw14/solo-037/internal/carriage"
    "github.com/wyw14/solo-037/internal/journal"
    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/track"
)

type Recovery struct { Journal *journal.Store; Track *track.State; Carriage carriage.RecoveryGate }

func (r Recovery) Restore(run model.RunID, fence uint64) bool {
    if !r.Track.Fenced(run, fence) { return false }
    return r.Carriage.Start(run, fence)
}

func (r Recovery) Steps(run model.RunID) []string {
    records := r.Journal.Records()
    decision := journal.Recover(records, run)
    return journal.RecoveryOrder(decision)
}
