package main

import (
	"log"
	"net/http"

	"github.com/zdub0is/adhd-productivity-app/api/internal/apikey"
	"github.com/zdub0is/adhd-productivity-app/api/internal/checkins"
	"github.com/zdub0is/adhd-productivity-app/api/internal/config"
	"github.com/zdub0is/adhd-productivity-app/api/internal/db"
	"github.com/zdub0is/adhd-productivity-app/api/internal/health"
	"github.com/zdub0is/adhd-productivity-app/api/internal/journal"
	"github.com/zdub0is/adhd-productivity-app/api/internal/planning"
	"github.com/zdub0is/adhd-productivity-app/api/internal/tasks"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	conn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer conn.Close()

	apiKeys := apikey.NewStore(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)

	// Bootstrapping new API keys requires the admin token, not an API key.
	mux.Handle("POST /api/v1/auth/keys",
		apikey.RequireAdminToken(cfg.AdminBootstrapToken)(http.HandlerFunc(apikey.NewHandler(apiKeys).Create)))

	v1 := http.NewServeMux()
	tasks.NewHandler(tasks.NewStore(conn)).Register(v1)
	planning.NewHandler(planning.NewStore(conn)).Register(v1)
	checkins.NewHandler(checkins.NewStore(conn)).Register(v1)
	journal.NewHandler(journal.NewStore(conn)).Register(v1)

	mux.Handle("/api/v1/", apikey.RequireAPIKey(apiKeys)(v1))

	log.Printf("listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
