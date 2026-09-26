package wave

import (
    "math"
    "time"

    "github.com/wyw14/solo-037/internal/model"
)

type Profile struct {
    Frequency float64 `json:"frequency"`
    Amplitude float64 `json:"amplitude"`
    Duration time.Duration `json:"duration"`
}

func NewProfile(frequency float64, amplitude float64, duration time.Duration) Profile {
    return Profile{Frequency: frequency, Amplitude: amplitude, Duration: duration}
}

func (p Profile) Sample(at time.Duration) float64 {
    if at < 0 || at > p.Duration {
        return 0
    }
    return p.Amplitude * math.Sin(2*math.Pi*p.Frequency*at.Seconds())
}

func (p Profile) Samples(step time.Duration) []float64 {
    result := make([]float64, 0)
    for elapsed := time.Duration(0); elapsed <= p.Duration; elapsed += step {
        result = append(result, p.Sample(elapsed))
    }
    return result
}

func ProfileEvidence(run model.RunID, revision model.RevisionID, p Profile, at time.Time) []model.Evidence {
    result := make([]model.Evidence, 0, 3)
    for index, value := range p.Samples(p.Duration / 3) {
        result = append(result, model.NewEvidence(run, revision, "paddle", model.EvidenceWaveTime, value, at.Add(time.Duration(index)*time.Millisecond)))
    }
    return result
}
