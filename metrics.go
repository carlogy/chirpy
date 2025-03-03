package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {

	loads := cfg.fileserverHits.Load()
	body := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", loads)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
