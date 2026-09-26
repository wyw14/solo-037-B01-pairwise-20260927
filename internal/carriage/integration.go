package carriage

import (
    "context"
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type Worker struct {
    mu sync.Mutex
    run model.RunID
    commands []model.DeviceCommand
    active bool
}

func NewWorker() *Worker { return &Worker{} }

func (w *Worker) Start(ctx context.Context, run model.RunID, command model.DeviceCommand) {
    w.mu.Lock()
    w.run = run
    w.commands = append(w.commands, command)
    w.active = true
    w.mu.Unlock()
    go func() {
        <-ctx.Done()
        w.mu.Lock()
        w.active = false
        w.mu.Unlock()
    }()
}

func (w *Worker) Abort() {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.active = false
}

func (w *Worker) Active() bool {
    w.mu.Lock()
    defer w.mu.Unlock()
    return w.active
}

func (w *Worker) Commands() []model.DeviceCommand {
    w.mu.Lock()
    defer w.mu.Unlock()
    result := make([]model.DeviceCommand, len(w.commands))
    copy(result, w.commands)
    return result
}
