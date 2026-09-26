package safety

import "sync"

type Interlocks struct { mu sync.Mutex; values map[string]bool }

func NewInterlocks() *Interlocks { return &Interlocks{values: map[string]bool{"track": true, "brake": true, "tank": true}} }
func (i *Interlocks) Set(name string, value bool) { i.mu.Lock(); defer i.mu.Unlock(); i.values[name] = value }
func (i *Interlocks) Ready() bool { i.mu.Lock(); defer i.mu.Unlock(); for _, value := range i.values { if !value { return false } }; return true }
func (i *Interlocks) Snapshot() map[string]bool { i.mu.Lock(); defer i.mu.Unlock(); result := make(map[string]bool, len(i.values)); for key, value := range i.values { result[key] = value }; return result }
