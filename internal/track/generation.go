package track

import (
    "sync/atomic"

    "github.com/wyw14/solo-037/internal/model"
)

type Generation struct {
    value atomic.Uint64
    revision atomic.Value
}

func NewGeneration() *Generation {
    generation := &Generation{}
    generation.revision.Store(model.NewRevisionID())
    return generation
}

func (g *Generation) Next() uint64 {
    return g.value.Add(1)
}

func (g *Generation) Revision() model.RevisionID {
    return g.revision.Load().(model.RevisionID)
}

func (g *Generation) Recalibrate() model.RevisionID {
    revision := model.NewRevisionID()
    g.revision.Store(revision)
    g.Next()
    return revision
}
