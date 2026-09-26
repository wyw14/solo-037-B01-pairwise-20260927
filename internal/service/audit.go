package service

import (
    "sort"
    "time"

    "github.com/wyw14/solo-037/internal/journal"
    "github.com/wyw14/solo-037/internal/model"
)

type AuditCheck struct {
    Name string `json:"name"`
    Passed bool `json:"passed"`
    Detail string `json:"detail"`
}

type AuditReport struct {
    Run model.RunID `json:"run"`
    Fence uint64 `json:"fence"`
    CheckedAt time.Time `json:"checked_at"`
    Passed bool `json:"passed"`
    Checks []AuditCheck `json:"checks"`
}

func AuditLaboratory(laboratory *Laboratory, run model.RunID, fence uint64, at time.Time) AuditReport {
    report := AuditReport{Run: run, Fence: fence, CheckedAt: at.UTC(), Passed: true}
    report.Add("identity", run != "", string(run))
    owner, currentFence, generation := laboratory.TrackState.Snapshot()
    report.Add("track_generation", generation > 0, unsignedDetail(generation))
    report.Add("track_owner_retired", owner == "", string(owner))
    report.Add("track_fence_retired", currentFence == 0, unsignedDetail(currentFence))
    occupancy := laboratory.Tank.Snapshot()
    report.Add("tank_section_released", len(occupancy) == 0, countDetail(len(occupancy)))
    samples := laboratory.Probes.Samples(run)
    report.Add("wave_samples_present", len(samples) > 0, countDetail(len(samples)))
    report.Add("wave_revision_current", sameRevision(samples), revisionDetail(samples))
    report.Add("wave_evidence_durable", durableEvidence(samples), durableDetail(samples))
    forces := laboratory.Forces.Values(run)
    report.Add("force_samples_present", len(forces) > 0, countDetail(len(forces)))
    report.Add("force_evidence_durable", durableEvidence(forces), durableDetail(forces))
    records := journal.RecordsForRun(laboratory.Store.Records(), run)
    report.Add("journal_records_present", len(records) > 0, countDetail(len(records)))
    recovery := journal.Recover(records, run)
    report.Add("recovery_fence_nonnegative", recovery.Fence <= fence, unsignedDetail(recovery.Fence))
    emergencyResult, emergencyFound := laboratory.Emergency.Result(run)
    report.Add("emergency_result_present", emergencyFound, boolDetail(emergencyFound))
    report.Add("emergency_secured", emergencyFound && emergencyResult.Secured, boolDetail(emergencyResult.Secured))
    report.Add("interlocks_ready", laboratory.Interlocks.Ready(), boolDetail(laboratory.Interlocks.Ready()))
    report.Sort()
    return report
}

func (r *AuditReport) Add(name string, passed bool, detail string) {
    r.Checks = append(r.Checks, AuditCheck{Name: name, Passed: passed, Detail: detail})
    r.Passed = r.Passed && passed
}

func (r *AuditReport) Sort() {
    sort.SliceStable(r.Checks, func(i int, j int) bool {
        return r.Checks[i].Name < r.Checks[j].Name
    })
}

func (r AuditReport) FailedChecks() []AuditCheck {
    result := make([]AuditCheck, 0)
    for _, check := range r.Checks {
        if !check.Passed {
            result = append(result, check)
        }
    }
    return result
}

func (r AuditReport) CheckMap() map[string]bool {
    result := make(map[string]bool, len(r.Checks))
    for _, check := range r.Checks {
        result[check.Name] = check.Passed
    }
    return result
}

func sameRevision(evidence []model.Evidence) bool {
    if len(evidence) == 0 {
        return false
    }
    revision := evidence[0].Revision
    for _, item := range evidence[1:] {
        if item.Revision != revision {
            return false
        }
    }
    return true
}

func durableEvidence(evidence []model.Evidence) bool {
    if len(evidence) == 0 {
        return false
    }
    for _, item := range evidence {
        if !item.Durable {
            return false
        }
    }
    return true
}

func revisionDetail(evidence []model.Evidence) string {
    if len(evidence) == 0 {
        return "none"
    }
    revisions := make(map[model.RevisionID]struct{})
    for _, item := range evidence {
        revisions[item.Revision] = struct{}{}
    }
    return countDetail(len(revisions))
}

func durableDetail(evidence []model.Evidence) string {
    count := 0
    for _, item := range evidence {
        if item.Durable {
            count++
        }
    }
    return countDetail(count) + "/" + countDetail(len(evidence))
}

func countDetail(value int) string {
    if value == 0 {
        return "0"
    }
    result := ""
    for value > 0 {
        result = string(rune('0'+value%10)) + result
        value /= 10
    }
    return result
}

func unsignedDetail(value uint64) string {
    if value == 0 {
        return "0"
    }
    result := ""
    for value > 0 {
        result = string(rune('0'+value%10)) + result
        value /= 10
    }
    return result
}

func boolDetail(value bool) string {
    if value {
        return "true"
    }
    return "false"
}
