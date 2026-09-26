package track

import (
    "fmt"
    "time"

    "github.com/wyw14/solo-037/internal/journal"
    "github.com/wyw14/solo-037/internal/model"
)

type Recovery struct {
    Coordinator *Coordinator
    State *State
}

func NewRecovery(coordinator *Coordinator, state *State) Recovery {
    return Recovery{Coordinator: coordinator, State: state}
}

func (r Recovery) Restore(run model.RunID, fence uint64) bool {
    if !r.State.Fenced(run, fence) {
        r.State.Assign(run, fence)
    }
    _, currentFence, _ := r.State.Snapshot()
    return currentFence == fence
}

func (r Recovery) Record(run model.RunID, fence uint64, store *journal.Store, at time.Time) error {
    record, err := journal.FenceRecord(run, fence, at)
    if err != nil {
        return err
    }
    return store.Append(record)
}

func LeaseDescription(lease Lease) string {
    return fmt.Sprintf("%s:%s:%d", lease.Section, lease.Run, lease.Fence)
}
