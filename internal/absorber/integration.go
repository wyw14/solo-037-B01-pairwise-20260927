package absorber

import (
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/probe"
)

func CalibrateFromProbe(controller *Controller, samples []model.Evidence, revision model.RevisionID) bool {
    if len(probe.SourceIDs(samples)) == 0 {
        return false
    }
    return controller.AddReflection(samples[0].Run, revision, len(samples)).Valid(revision)
}

func ActivateAfterBoundary(controller *Controller, at time.Time) error {
    if err := controller.RequirePassiveGate(); err != nil {
        return err
    }
    return controller.SwitchMode("active", at)
}
