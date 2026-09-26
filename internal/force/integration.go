package force

import (
    "time"

    "github.com/wyw14/solo-037/internal/journal"
    "github.com/wyw14/solo-037/internal/model"
)

func Persist(receiver *Receiver, store *journal.Store, run model.RunID, revision model.RevisionID, value float64, at time.Time) error {
    evidence := receiver.Record(run, revision, value, at)
    record, err := journal.EvidenceRecord(evidence, at)
    if err != nil {
        return err
    }
    return store.Append(record)
}

func Complete(receiver *Receiver, run model.RunID) bool { return receiver.Durable(run) }
