package tank

import (
    "errors"
    "sync"
    "time"
)

type BoundaryState struct {
    Mode string `json:"mode"`
    Angle float64 `json:"angle"`
    Isolated bool `json:"isolated"`
    ConfirmedAt time.Time `json:"confirmed_at"`
}

type Boundary struct {
    mu sync.Mutex
    state BoundaryState
    events []string
}

func NewBoundary() *Boundary {
    return &Boundary{state: BoundaryState{Mode: "passive", Angle: 35}}
}

func (b *Boundary) IsolatePassive(at time.Time) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.state.Mode != "passive" {
        return errors.New("passive boundary is not active")
    }
    b.state.Angle = 0
    b.state.Isolated = true
    b.state.ConfirmedAt = at.UTC()
    b.events = append(b.events, "passive_isolated")
    return nil
}

func (b *Boundary) Activate(mode string, at time.Time) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    if !b.state.Isolated {
        return errors.New("old boundary not isolated")
    }
    b.state.Mode = mode
    b.state.ConfirmedAt = at.UTC()
    b.events = append(b.events, mode+"_active")
    return nil
}

func (b *Boundary) State() BoundaryState {
    b.mu.Lock()
    defer b.mu.Unlock()
    return b.state
}

func (b *Boundary) Events() []string {
    b.mu.Lock()
    defer b.mu.Unlock()
    result := make([]string, len(b.events))
    copy(result, b.events)
    return result
}

func (b *Boundary) ResetPassive() {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.state = BoundaryState{Mode: "passive", Angle: 35}
    b.events = nil
}
