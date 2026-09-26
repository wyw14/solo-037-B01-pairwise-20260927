package absorber

import (
    "errors"
    "sync"
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/tank"
    "github.com/wyw14/solo-037/internal/wave"
)

type Controller struct {
    mu sync.Mutex
    boundary *tank.Boundary
    mode string
    windows map[model.RunID]wave.ReflectionWindow
}

func NewController(boundary *tank.Boundary) *Controller { return &Controller{boundary: boundary, mode: "passive", windows: make(map[model.RunID]wave.ReflectionWindow)} }

func (c *Controller) SwitchMode(mode string, at time.Time) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    if mode == c.mode {
        return nil
    }
    if err := c.boundary.IsolatePassive(at); err != nil {
        return err
    }
    if err := c.boundary.Activate(mode, at); err != nil {
        return err
    }
    c.mode = mode
    return nil
}

func (c *Controller) Mode() string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.mode
}

func (c *Controller) AddReflection(run model.RunID, revision model.RevisionID, samples int) wave.ReflectionWindow {
    c.mu.Lock()
    defer c.mu.Unlock()
    window, ok := c.windows[run]
    if !ok {
        window = wave.NewReflectionWindow(run, revision)
    }
    window = window.Add(revision, samples)
    c.windows[run] = window
    return window
}

func (c *Controller) ReflectionReady(run model.RunID, revision model.RevisionID) bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    window, ok := c.windows[run]
    return ok && window.Valid(revision)
}

func (c *Controller) RequirePassiveGate() error {
    if !c.boundary.State().Isolated {
        return errors.New("passive gate pending")
    }
    return nil
}
