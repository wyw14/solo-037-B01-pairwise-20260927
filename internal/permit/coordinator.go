package permit

import (
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type Coordinator struct {
    mu sync.Mutex
    evaluator *Evaluator
    states map[model.RunID]string
}

func NewCoordinator(evaluator *Evaluator) *Coordinator { return &Coordinator{evaluator: evaluator, states: make(map[model.RunID]string)} }

func (c *Coordinator) Update(run model.RunID) Result {
    c.mu.Lock()
    defer c.mu.Unlock()
    result := c.evaluator.Evaluate(run)
    if result.Allowed {
        c.states[run] = "ready"
    } else {
        c.states[run] = "waiting"
    }
    return result
}

func (c *Coordinator) State(run model.RunID) string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.states[run]
}

func (c *Coordinator) Evaluator() *Evaluator { return c.evaluator }
