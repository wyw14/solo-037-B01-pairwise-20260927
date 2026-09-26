package permit

import "github.com/wyw14/solo-037/internal/model"

type Recovery struct {
    Stage model.RunStage `json:"stage"`
    Permit Result `json:"permit"`
}

func Restore(stage model.RunStage, result Result) Recovery {
    if stage == model.StageComplete && !result.Allowed {
        stage = model.StageSettling
    }
    return Recovery{Stage: stage, Permit: result}
}

func (r Recovery) CanClear() bool { return r.Stage == model.StageComplete && r.Permit.Allowed }

func (r Recovery) Description() string {
    if r.CanClear() {
        return "clearance permitted"
    }
    return "clearance held"
}
