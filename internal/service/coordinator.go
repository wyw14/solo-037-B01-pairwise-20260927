package service

import (
    "sync"
    "time"

    "github.com/wyw14/solo-037/internal/model"
)

type Coordinator struct { Runtime *Runtime; mu sync.Mutex; generation uint64 }

func NewCoordinator(runtime *Runtime) *Coordinator { return &Coordinator{Runtime: runtime} }

func (c *Coordinator) Start(operation model.Operation) (model.Operation, error) {
    c.mu.Lock()
    c.generation++
    c.mu.Unlock()
    return c.Runtime.Create(operation)
}

func (c *Coordinator) Abort(run model.RunID) error { return c.Runtime.Abort(run, time.Now()) }

func (c *Coordinator) Generation() uint64 { c.mu.Lock(); defer c.mu.Unlock(); return c.generation }
