package journal

import "github.com/wyw14/solo-037/internal/model"

type RecoveryDecision struct {
    Run model.RunID `json:"run"`
    Stage model.RunStage `json:"stage"`
    Fence uint64 `json:"fence"`
    EvidenceReady bool `json:"evidence_ready"`
    CanClear bool `json:"can_clear"`
}

func Recover(records []Record, run model.RunID) RecoveryDecision {
    replay := Replay(records, run)
    evidenceReady := replay.Evidence >= 2
    stage := replay.Stage
    if stage == model.StageComplete && !evidenceReady {
        stage = model.StageSettling
    }
    return RecoveryDecision{
        Run: run,
        Stage: stage,
        Fence: replay.Fence,
        EvidenceReady: evidenceReady,
        CanClear: stage == model.StageComplete && evidenceReady,
    }
}

func RecoveryOrder(decision RecoveryDecision) []string {
    steps := []string{"journal_loaded"}
    if decision.Fence > 0 {
        steps = append(steps, "track_fence_loaded")
    }
    if decision.EvidenceReady {
        steps = append(steps, "measurement_evidence_loaded")
    }
    steps = append(steps, "controllers_ready")
    return steps
}

func MergeRecovery(left RecoveryDecision, right RecoveryDecision) RecoveryDecision {
    if right.Fence > left.Fence {
        left.Fence = right.Fence
    }
    left.EvidenceReady = left.EvidenceReady && right.EvidenceReady
    left.CanClear = left.CanClear && right.CanClear
    if model.StageIndex(right.Stage) < model.StageIndex(left.Stage) {
        left.Stage = right.Stage
    }
    return left
}
