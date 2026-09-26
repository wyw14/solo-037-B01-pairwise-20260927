package telemetry

import (
    "sync"
    "time"
)

type Span struct { Name string; Started time.Time; Ended time.Time; Attributes map[string]string }
type Tracer struct { mu sync.Mutex; spans []Span }

func NewTracer() *Tracer { return &Tracer{} }
func (t *Tracer) Start(name string, at time.Time) int { t.mu.Lock(); defer t.mu.Unlock(); t.spans = append(t.spans, Span{Name: name, Started: at, Attributes: map[string]string{}}); return len(t.spans)-1 }
func (t *Tracer) End(index int, at time.Time) { t.mu.Lock(); defer t.mu.Unlock(); if index >= 0 && index < len(t.spans) { t.spans[index].Ended = at } }
func (t *Tracer) Spans() []Span { t.mu.Lock(); defer t.mu.Unlock(); return append([]Span{}, t.spans...) }
