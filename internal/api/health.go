package api

import (
    "net/http"
    "time"
)

type Health struct { StartedAt time.Time; Version string }

func (h Health) Status() map[string]any {
    return map[string]any{"started_at": h.StartedAt, "version": h.Version, "uptime": time.Since(h.StartedAt).String()}
}

func (h Health) Handler(w http.ResponseWriter, _ *http.Request) {
    writeJSON(w, h.Status())
}
