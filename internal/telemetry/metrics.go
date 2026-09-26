package telemetry

import "sync/atomic"

type Metrics struct { requests atomic.Uint64; events atomic.Uint64; faults atomic.Uint64 }

func (m *Metrics) Request() { m.requests.Add(1) }
func (m *Metrics) Event() { m.events.Add(1) }
func (m *Metrics) Fault() { m.faults.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 { return map[string]uint64{"requests": m.requests.Load(), "events": m.events.Load(), "faults": m.faults.Load()} }
