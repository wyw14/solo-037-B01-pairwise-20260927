package permit

import (
    "github.com/wyw14/solo-037/internal/carriage"
    "github.com/wyw14/solo-037/internal/model"
)

func BindBrake(evaluator *Evaluator, controller *carriage.Controller, run model.RunID, commandID string, epoch uint64) bool {
    if !controller.Stopped(run, commandID, epoch) {
        return false
    }
    receipt, ok := controllerReceipt(controller, run)
    if !ok {
        return false
    }
    evaluator.ObserveBrake(run, receipt)
    return evaluator.AcceptBrake(run, commandID, epoch)
}

func controllerReceipt(controller *carriage.Controller, run model.RunID) (model.DeviceReceipt, bool) {
    command := controller.Commands(run)
    if command.ID == "" {
        return model.DeviceReceipt{}, false
    }
    return model.DeviceReceipt{CommandID: command.ID, Epoch: command.Epoch, Velocity: 0, Confirmed: controller.Stopped(run, command.ID, command.Epoch)}, true
}

func SyncWaveSources(evaluator *Evaluator, run model.RunID, sources int, calm bool) Result {
    evaluator.SetSources(run, sources)
    evaluator.SetCalm(run, calm)
    return evaluator.Evaluate(run)
}
