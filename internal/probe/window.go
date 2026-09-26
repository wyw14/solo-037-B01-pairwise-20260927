package probe

import (
    "sync"
    "time"

    "github.com/wyw14/solo-037/internal/model"
)

type Sample struct {
    Run model.RunID `json:"run"`
    Revision model.RevisionID `json:"revision"`
    Source string `json:"source"`
    Height float64 `json:"height"`
    At time.Time `json:"at"`
}

type CalmWindow struct {
    mu sync.Mutex
    run model.RunID
    revision model.RevisionID
    last time.Time
    started time.Time
    required time.Duration
    valid bool
}

func NewCalmWindow(run model.RunID, revision model.RevisionID, required time.Duration) *CalmWindow {
    return &CalmWindow{run: run, revision: revision, required: required}
}

func (w *CalmWindow) Observe(sample Sample) {
    w.mu.Lock()
    defer w.mu.Unlock()
    if sample.Run != w.run || sample.Revision != w.revision || sample.Height > 0.02 {
        w.started = time.Time{}
        w.valid = false
        w.last = sample.At
        return
    }
    if w.last.IsZero() || sample.At.Sub(w.last) > 1500*time.Millisecond {
        w.started = sample.At
    }
    if w.started.IsZero() {
        w.started = sample.At
    }
    w.last = sample.At
    w.valid = sample.At.Sub(w.started) >= w.required
}

func (w *CalmWindow) Valid() bool {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.valid
}

func (w *CalmWindow) Reset(revision model.RevisionID) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.revision = revision
    w.last = time.Time{}
    w.started = time.Time{}
    w.valid = false
}
