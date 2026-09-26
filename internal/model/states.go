package model

import (
    "errors"
    "fmt"
)

type RunStage string

const (
    StageConfigured RunStage = "configured"
    StageArming RunStage = "arming"
    StageAccelerating RunStage = "accelerating"
    StageMeasuring RunStage = "measuring"
    StageDecelerating RunStage = "decelerating"
    StageSettling RunStage = "settling"
    StageComplete RunStage = "complete"
    StageAborted RunStage = "aborted"
)

var stageOrder = []RunStage{
    StageConfigured,
    StageArming,
    StageAccelerating,
    StageMeasuring,
    StageDecelerating,
    StageSettling,
    StageComplete,
}

func StageIndex(stage RunStage) int {
    for index, candidate := range stageOrder {
        if candidate == stage {
            return index
        }
    }
    return -1
}

func CanAdvance(from RunStage, to RunStage) bool {
    if to == StageAborted && from != StageComplete {
        return true
    }
    fromIndex := StageIndex(from)
    toIndex := StageIndex(to)
    return fromIndex >= 0 && toIndex == fromIndex+1
}

func Advance(from RunStage, to RunStage) (RunStage, error) {
    if !CanAdvance(from, to) {
        return from, fmt.Errorf("invalid stage transition %s to %s", from, to)
    }
    return to, nil
}

func Terminal(stage RunStage) bool {
    return stage == StageComplete || stage == StageAborted
}

func ValidateStage(stage RunStage) error {
    if StageIndex(stage) >= 0 || stage == StageAborted {
        return nil
    }
    return errors.New("unknown run stage")
}

func StageSequence() []RunStage {
    result := make([]RunStage, len(stageOrder))
    copy(result, stageOrder)
    return result
}
