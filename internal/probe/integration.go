package probe

import (
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/tank"
)

func RecordToSection(coordinator *Coordinator, sections *tank.Sections, section string, run model.RunID, revision model.RevisionID, source string, height float64, at time.Time) model.Evidence {
    evidence := coordinator.Record(run, revision, source, height, at)
    sections.RecordWave(section, height, at)
    return evidence
}

func ResetAfterGap(coordinator *Coordinator, revision model.RevisionID) {
    coordinator.Window.Reset(revision)
}

func SourceIDs(samples []model.Evidence) []string {
    values := make([]string, 0, len(samples))
    seen := make(map[string]struct{})
    for _, sample := range samples {
        if _, ok := seen[sample.SourceID]; !ok {
            values = append(values, sample.SourceID)
            seen[sample.SourceID] = struct{}{}
        }
    }
    return values
}
