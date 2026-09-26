package carriage

import (
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type PositionWindow struct {
    Revision model.RevisionID `json:"revision"`
    Start float64 `json:"start"`
    End float64 `json:"end"`
    Entered bool `json:"entered"`
}

type PositionState struct {
    mu sync.Mutex
    revision model.RevisionID
    position float64
    window PositionWindow
}

func NewPositionState(revision model.RevisionID) *PositionState {
    return &PositionState{revision: revision}
}

func (s *PositionState) UpdateRevision(revision model.RevisionID) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.revision = revision
    s.window = PositionWindow{Revision: revision}
}

func (s *PositionState) SetWindow(start float64, end float64) PositionWindow {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.window = PositionWindow{Revision: s.revision, Start: start, End: end}
    return s.window
}

func (s *PositionState) Observe(position float64) PositionWindow {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.position = position
    if position >= s.window.Start && position <= s.window.End && s.window.Revision == s.revision {
        s.window.Entered = true
    }
    return s.window
}

func (s *PositionState) Snapshot() PositionWindow {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.window
}

func (s *PositionState) Position() float64 {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.position
}
