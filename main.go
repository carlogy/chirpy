package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/carlogy/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
	polka_key      string
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	s := os.Getenv("SECRET")
	pk := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	log.Print(dbQueries)

	cfg := &apiConfig{}
	cfg.dbQueries = dbQueries
	cfg.platform = platform
	cfg.secret = s
	cfg.polka_key = pk

	const port = "8080"
	const rootDir = http.Dir(".")

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(rootDir))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("POST /api/validate_chirp", HandlerBody)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUsers)

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)

	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)

	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeletechirp)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)

	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUsers)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Printf("Access: http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
