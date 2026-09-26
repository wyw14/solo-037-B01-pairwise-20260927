package tank

import (
    "errors"
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type Occupancy struct {
    Section string `json:"section"`
    Run model.RunID `json:"run"`
    Carriage string `json:"carriage"`
    Lease model.LeaseID `json:"lease"`
    Fence uint64 `json:"fence"`
}

type Registry struct {
    mu sync.Mutex
    nextFence uint64
    occupancy map[string]Occupancy
    depthRevision model.RevisionID
    depthMeters float64
}

func NewRegistry(depth float64) *Registry {
    return &Registry{
        occupancy: make(map[string]Occupancy),
        depthRevision: model.NewRevisionID(),
        depthMeters: depth,
    }
}

func (r *Registry) Claim(section string, run model.RunID, carriage string) (Occupancy, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    if current, exists := r.occupancy[section]; exists && current.Run != run {
        return Occupancy{}, errors.New("section occupied")
    }
    r.nextFence++
    occupancy := Occupancy{
        Section: section,
        Run: run,
        Carriage: carriage,
        Lease: model.NewLeaseID(),
        Fence: r.nextFence,
    }
    r.occupancy[section] = occupancy
    return occupancy, nil
}

func (r *Registry) ObserveFree(section string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    _, exists := r.occupancy[section]
    return !exists
}

func (r *Registry) CommitObserved(section string, run model.RunID, carriage string) (Occupancy, error) {
    return r.Claim(section, run, carriage)
}

func (r *Registry) Release(section string, run model.RunID, fence uint64) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    current, exists := r.occupancy[section]
    if !exists || current.Run != run || current.Fence != fence {
        return false
    }
    delete(r.occupancy, section)
    return true
}

func (r *Registry) Occupancy(section string) (Occupancy, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    current, exists := r.occupancy[section]
    return current, exists
}

func (r *Registry) Snapshot() []Occupancy {
    r.mu.Lock()
    defer r.mu.Unlock()
    result := make([]Occupancy, 0, len(r.occupancy))
    for _, occupancy := range r.occupancy {
        result = append(result, occupancy)
    }
    return result
}
