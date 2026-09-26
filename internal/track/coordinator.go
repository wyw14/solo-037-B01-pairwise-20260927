package track

import (
    "errors"
    "sync"
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/tank"
)

type Lease struct {
    Section string `json:"section"`
    Run model.RunID `json:"run"`
    Carriage string `json:"carriage"`
    Fence uint64 `json:"fence"`
    LeaseID model.LeaseID `json:"lease_id"`
    AcquiredAt time.Time `json:"acquired_at"`
}

type Coordinator struct {
    mu sync.Mutex
    registry *tank.Registry
    leases map[model.RunID]Lease
}

func NewCoordinator(registry *tank.Registry) *Coordinator {
    return &Coordinator{registry: registry, leases: make(map[model.RunID]Lease)}
}

func (c *Coordinator) Arm(run model.RunID, section string, carriage string, at time.Time) (Lease, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    occupancy, err := c.registry.Claim(section, run, carriage)
    if err != nil {
        return Lease{}, err
    }
    lease := Lease{Section: section, Run: run, Carriage: carriage, Fence: occupancy.Fence, LeaseID: occupancy.Lease, AcquiredAt: at.UTC()}
    c.leases[run] = lease
    return lease, nil
}

func (c *Coordinator) Release(run model.RunID) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    lease, ok := c.leases[run]
    if !ok {
        return errors.New("lease not found")
    }
    if !c.registry.Release(lease.Section, run, lease.Fence) {
        return errors.New("lease fence mismatch")
    }
    delete(c.leases, run)
    return nil
}

func (c *Coordinator) ReleaseWithFence(run model.RunID, fence uint64) bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    lease, ok := c.leases[run]
    if !ok || lease.Fence != fence {
        return false
    }
    if !c.registry.Release(lease.Section, run, fence) {
        return false
    }
    delete(c.leases, run)
    return true
}

func (c *Coordinator) Lease(run model.RunID) (Lease, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    lease, ok := c.leases[run]
    return lease, ok
}

func (c *Coordinator) Active() []Lease {
    c.mu.Lock()
    defer c.mu.Unlock()
    result := make([]Lease, 0, len(c.leases))
    for _, lease := range c.leases {
        result = append(result, lease)
    }
    return result
}
