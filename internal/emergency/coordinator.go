package emergency

import (
    "time"

    "github.com/wyw14/solo-037/internal/model"
)

type Coordinator struct { Brake *Brake }

func NewCoordinator(brake *Brake) *Coordinator { return &Coordinator{Brake: brake} }

func (c *Coordinator) Stop(run model.RunID, commandID string, epoch uint64, failure error, at time.Time) Result {
    return c.Brake.Execute(run, commandID, epoch, failure, at)
}

func (c *Coordinator) Secured(run model.RunID) bool {
    result, ok := c.Brake.Result(run)
    return ok && result.Secured
}
