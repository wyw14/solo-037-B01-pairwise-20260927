package track

import (
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type State struct {
    mu sync.Mutex
    generation uint64
    owner model.RunID
    fence uint64
}

func (s *State) Assign(run model.RunID, fence uint64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.generation++
    s.owner = run
    s.fence = fence
}

func (s *State) Fenced(run model.RunID, fence uint64) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.owner == run && s.fence == fence
}

func (s *State) Snapshot() (model.RunID, uint64, uint64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.owner, s.fence, s.generation
}

func (s *State) Clear(run model.RunID, fence uint64) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.owner != run || s.fence != fence {
        return false
    }
    s.owner = ""
    s.fence = 0
    return true
}
