package main

import (
	"flag"
	"log/slog"

	"proudtrails.org/trails/internal"
)

func main() {
	generateOpenAPI := flag.Bool("generate-openapi", false, "Generate OpenAPI spec and write to ../generated/openapi.json")
	flag.Parse()

	log := slog.With("app", "trailsapp")

	if *generateOpenAPI {
		if err := internal.GenerateSpec(log); err != nil {
			log.Error("Failed to generate OpenAPI spec", "error", err)
		}
		return
	}

	if err := internal.App(log); err != nil {
		log.Error("Failed starting app", "error", err)
	}

	log.Info("Shutdown")
}
