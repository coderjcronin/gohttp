package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/coderjcronin/gohttp/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	secret         string
	expires        time.Duration
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	apiSecret := os.Getenv("SECRET")
	stringExp := os.Getenv("EXPIRES")

	durationExp, _ := time.ParseDuration(stringExp)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error loading database: %s", err)
		return
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		secret:         apiSecret,
		expires:        durationExp,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/chirps", apiCfg.apiGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.apiGetChirp)

	mux.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	mux.HandleFunc("POST /api/users", apiCfg.apiAddUser)
	mux.HandleFunc("POST /api/login", apiCfg.apiCheckLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.apiRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.apiRevokeToken)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.apiPolkaWebhook)

	mux.HandleFunc("PUT /api/users", apiCfg.apidUpdateUser)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.apiDeleteChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
