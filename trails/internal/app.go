package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func setupAPI(r chi.Router) huma.API {
	api := humachi.New(r, huma.DefaultConfig("Proud Trails API", "1.0.0"))
	registerTrailRoutes(api)
	return api
}

func App(log *slog.Logger) error {
	log.Info("Starting app")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	setupAPI(r)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Error("Failed to start HTTP server", "error", err)
		return err
	}

	log.Info("Shutting down app")
	return nil
}

func GenerateSpec(log *slog.Logger) error {
	r := chi.NewRouter()
	api := setupAPI(r)

	b, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling OpenAPI spec: %w", err)
	}

	// Relative to proudlands/trails/ — resolves to proudlands/generated/
	outDir := filepath.Join("..", "generated")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	outPath := filepath.Join(outDir, "openapi.json")
	if err := os.WriteFile(outPath, b, 0644); err != nil {
		return fmt.Errorf("writing OpenAPI spec: %w", err)
	}

	log.Info("OpenAPI spec written", "path", outPath)
	return nil
}
