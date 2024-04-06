package main

import (
	"avito/assignment/config"
	"avito/assignment/internal/server"
	"avito/assignment/pkg/constant"
	"avito/assignment/pkg/traces"
	"context"
	"github.com/jmoiron/sqlx"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	traces.NewTracer(ctx, *cfg)
	defer func(tracer *tracesdk.TracerProvider) {
		if err := tracer.Shutdown(ctx); err != nil {
			log.Printf(err.Error())
		} else {
			log.Println("Jaeger closed properly")
		}
	}(constant.Tracer)

	psqlDB, err := server.NewDB(*cfg)
	if err != nil {
		log.Fatalf("psqlDB: %v", err)
	}
	defer func(psqlDB *sqlx.DB) {
		if err := psqlDB.Close(); err != nil {
			log.Printf(err.Error())
		} else {
			log.Println("PostgresSQL closed properly")
		}
	}(psqlDB)

	redis, err := server.NewRedisClient(*cfg)
	if err != nil {
		log.Printf("Failed to connect to redis: %s", err.Error())
	}

	s := server.NewServer(
		cfg,
		psqlDB,
		redis,
	)
	if err = s.Run(); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
