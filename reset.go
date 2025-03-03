package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {


	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}

	// load := cfg.fileserverHits.Load()
	// swapped := cfg.fileserverHits.CompareAndSwap(load, 0)
	// if !swapped {
	// 	log.Println("failed to swap")
	// }
	// loads := strconv.Itoa(int(cfg.fileserverHits.Load()))

	// w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Hits: " + loads))
	//
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`Error "Resetting" users table`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Reset users table"))
}
