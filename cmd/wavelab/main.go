package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wyw14/solo-037/internal/api"
	"github.com/wyw14/solo-037/internal/service"
)

func main() {
	dataDirectory := os.Getenv("WAVELAB_DATA_DIR")
	if dataDirectory == "" {
		dataDirectory = filepath.Join("data", "wavelab")
	}
	runtime, err := service.OpenRuntime(dataDirectory)
	if err != nil {
		log.Fatal(err)
	}
	router := api.NewRouter(runtime)
	mux := http.NewServeMux()
	mux.Handle("/", router.Handler())
	mux.HandleFunc("/api/snapshot", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(runtime.Snapshot())
	})
	address := os.Getenv("WAVELAB_ADDR")
	if address == "" {
		address = "127.0.0.1:21244"
	}
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("wavelab listening on %s", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
