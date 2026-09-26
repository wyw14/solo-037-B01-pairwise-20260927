package carriage

import (
    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/track"
)

type RecoveryGate struct {
    TrackState *track.State
    Controller *Controller
}

func NewRecoveryGate(state *track.State, controller *Controller) RecoveryGate {
    return RecoveryGate{TrackState: state, Controller: controller}
}

func (g RecoveryGate) Start(run model.RunID, fence uint64) bool {
    return g.TrackState.Fenced(run, fence)
}

func (g RecoveryGate) CommandAllowed(run model.RunID, fence uint64) bool {
    return g.Start(run, fence)
}

func (g RecoveryGate) Stop(run model.RunID) bool {
    return g.Controller.Velocity(run) == 0
}
