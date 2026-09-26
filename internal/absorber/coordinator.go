package absorber

import (
    "time"

    "github.com/wyw14/solo-037/internal/model"
)

type Coordinator struct { Controller *Controller }

func NewCoordinator(controller *Controller) *Coordinator { return &Coordinator{Controller: controller} }

func (c *Coordinator) ConfigureActive(at time.Time) error { return c.Controller.SwitchMode("active", at) }

func (c *Coordinator) ObserveReflection(run model.RunID, revision model.RevisionID, samples int) bool {
    return c.Controller.AddReflection(run, revision, samples).Valid(revision)
}

func (c *Coordinator) Ready(run model.RunID, revision model.RevisionID) bool {
    return c.Controller.ReflectionReady(run, revision)
}
