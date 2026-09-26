package permit

import (
    "sync"

    "github.com/wyw14/solo-037/internal/model"
)

type Evaluator struct {
    mu sync.Mutex
    brake map[model.RunID]model.DeviceReceipt
    sourceCount map[model.RunID]int
    calm map[model.RunID]bool
}

func NewEvaluator() *Evaluator { return &Evaluator{brake: make(map[model.RunID]model.DeviceReceipt), sourceCount: make(map[model.RunID]int), calm: make(map[model.RunID]bool)} }

func (e *Evaluator) AcceptBrake(run model.RunID, commandID string, epoch uint64) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    receipt, ok := e.brake[run]
    return ok && receipt.CommandID == commandID && receipt.Epoch == epoch && receipt.Confirmed && receipt.Velocity == 0
}

func (e *Evaluator) ObserveBrake(run model.RunID, receipt model.DeviceReceipt) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.brake[run] = receipt
}

func (e *Evaluator) SetSources(run model.RunID, count int) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.sourceCount[run] = count
}

func (e *Evaluator) SetCalm(run model.RunID, value bool) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.calm[run] = value
}

func (e *Evaluator) WavePermit(run model.RunID) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.sourceCount[run] >= 2 && e.calm[run]
}

func (e *Evaluator) SourceCount(run model.RunID) int {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.sourceCount[run]
}

type Result struct {
    Run model.RunID `json:"run"`
    Allowed bool `json:"allowed"`
    Reason string `json:"reason"`
}

func (e *Evaluator) Evaluate(run model.RunID) Result {
    if e.WavePermit(run) {
        return Result{Run: run, Allowed: true, Reason: "independent evidence and calm window"}
    }
    return Result{Run: run, Allowed: false, Reason: "evidence prerequisites incomplete"}
}
