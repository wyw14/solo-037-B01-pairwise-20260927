package permit

import (
    "sort"
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/probe"
)

type Requirement struct {
    Name string `json:"name"`
    Satisfied bool `json:"satisfied"`
    Detail string `json:"detail"`
}

type Decision struct {
    Run model.RunID `json:"run"`
    Allowed bool `json:"allowed"`
    Requirements []Requirement `json:"requirements"`
    EvaluatedAt time.Time `json:"evaluated_at"`
}

type Policy struct {
    MinimumSources int
    MaximumGap time.Duration
    MaximumWaveRMS float64
}

func DefaultPolicy() Policy {
    return Policy{MinimumSources: 2, MaximumGap: 1500 * time.Millisecond, MaximumWaveRMS: 0.03}
}

func (p Policy) Evaluate(run model.RunID, samples []model.Evidence, calm bool, brake bool, at time.Time) Decision {
    matching := probe.FilterRun(samples, run)
    analysis := probe.Analyze(matching)
    requirements := []Requirement{
        {Name: "independent_sources", Satisfied: len(analysis.Sources) >= p.MinimumSources, Detail: sourceDetail(analysis.Sources)},
        {Name: "continuous_window", Satisfied: probe.Continuous(matching, p.MaximumGap), Detail: analysis.Span.String()},
        {Name: "durable_wave_evidence", Satisfied: probe.Durable(matching), Detail: durabilityDetail(matching)},
        {Name: "calm_water", Satisfied: calm && analysis.RMS <= p.MaximumWaveRMS, Detail: waveDetail(analysis.RMS)},
        {Name: "carriage_brake", Satisfied: brake, Detail: brakeDetail(brake)},
    }
    allowed := true
    for _, requirement := range requirements {
        allowed = allowed && requirement.Satisfied
    }
    return Decision{Run: run, Allowed: allowed, Requirements: requirements, EvaluatedAt: at.UTC()}
}

func sourceDetail(sources []string) string {
    values := append([]string{}, sources...)
    sort.Strings(values)
    result := ""
    for index, value := range values {
        if index > 0 {
            result += ","
        }
        result += value
    }
    return result
}

func durabilityDetail(samples []model.Evidence) string {
    durable := 0
    for _, sample := range samples {
        if sample.Durable {
            durable++
        }
    }
    return formatCount(durable, len(samples))
}

func waveDetail(value float64) string {
    if value <= 0.01 {
        return "calm"
    }
    if value <= 0.03 {
        return "settling"
    }
    return "active"
}

func brakeDetail(value bool) string {
    if value {
        return "current command stopped"
    }
    return "current command proof missing"
}

func formatCount(current int, total int) string {
    digits := func(value int) string {
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
    return digits(current) + "/" + digits(total)
}

func Missing(decision Decision) []string {
    result := make([]string, 0)
    for _, requirement := range decision.Requirements {
        if !requirement.Satisfied {
            result = append(result, requirement.Name)
        }
    }
    sort.Strings(result)
    return result
}

func RequirementMap(decision Decision) map[string]bool {
    result := make(map[string]bool, len(decision.Requirements))
    for _, requirement := range decision.Requirements {
        result[requirement.Name] = requirement.Satisfied
    }
    return result
}
