package emergency

import (
    "errors"
    "sync"
    "time"

    "github.com/wyw14/solo-037/internal/carriage"
    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/wave"
)

type Result struct {
    Run model.RunID `json:"run"`
    CarriageError string `json:"carriage_error"`
    WaveStopped bool `json:"wave_stopped"`
    Secured bool `json:"secured"`
}

type Brake struct {
    mu sync.Mutex
    carriage *carriage.Controller
    wave *wave.Controller
    results map[model.RunID]Result
}

func NewBrake(carriageController *carriage.Controller, waveController *wave.Controller) *Brake {
    return &Brake{carriage: carriageController, wave: waveController, results: make(map[model.RunID]Result)}
}

func (b *Brake) Execute(run model.RunID, commandID string, epoch uint64, carriageFailure error, at time.Time) Result {
    b.mu.Lock()
    defer b.mu.Unlock()
    result := Result{Run: run}
    if carriageFailure != nil {
        result.CarriageError = carriageFailure.Error()
    } else {
        _, err := b.carriage.Brake(run, commandID, epoch, at)
        if err != nil {
            result.CarriageError = err.Error()
        }
    }
    b.wave.Stop()
    result.WaveStopped = !b.wave.Active()
    result.Secured = result.CarriageError == "" && result.WaveStopped
    b.results[run] = result
    return result
}

func (b *Brake) Result(run model.RunID) (Result, bool) {
    b.mu.Lock()
    defer b.mu.Unlock()
    result, ok := b.results[run]
    return result, ok
}

func (r Result) Error() error {
    if r.CarriageError != "" {
        return errors.New(r.CarriageError)
    }
    if !r.Secured {
        return errors.New("emergency not secured")
    }
    return nil
}
