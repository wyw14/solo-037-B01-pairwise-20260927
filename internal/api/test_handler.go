package api

import (
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/wyw14/solo-037/internal/model"
    "github.com/wyw14/solo-037/internal/service"
)

type TestHandler struct { Service TestService }

type TestService interface {
    Create(model.Operation) (model.Operation, error)
    Abort(model.RunID, time.Time) error
}

func (h TestHandler) Create(w http.ResponseWriter, request *http.Request) {
    operation := model.NewOperation("water-run", "model", "section-a", "carriage-1", time.Now())
    result, err := h.Service.Create(operation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }
    _ = json.NewEncoder(w).Encode(result)
}

func (h TestHandler) Abort(w http.ResponseWriter, request *http.Request) {
    err := h.Service.Abort(model.RunID(request.URL.Query().Get("run")), time.Now())
    if errors.Is(err, service.ErrRunNotFound) {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
