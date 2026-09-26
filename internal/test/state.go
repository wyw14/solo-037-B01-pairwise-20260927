package test

import (
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type Registry struct {
    mu sync.Mutex
    operations map[model.RunID]model.Operation
}

func NewRegistry() *Registry { return &Registry{operations: make(map[model.RunID]model.Operation)} }

func (r *Registry) Put(operation model.Operation) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.operations[operation.ID] = operation
}

func (r *Registry) Get(run model.RunID) (model.Operation, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    operation, ok := r.operations[run]
    return operation, ok
}

func (r *Registry) Active() []model.Operation {
    r.mu.Lock()
    defer r.mu.Unlock()
    result := make([]model.Operation, 0, len(r.operations))
    for _, operation := range r.operations {
        if operation.Active() {
            result = append(result, operation)
        }
    }
    return result
}
