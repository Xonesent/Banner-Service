package server

import (
	"avito/assignment/config"
	"avito/assignment/pkg/error_handler"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg   *config.Config
	pgDB  *sqlx.DB
	fiber *fiber.App
	redis *redis.Client
}

func NewServer(
	cfg *config.Config,
	pgDB *sqlx.DB,
	redis *redis.Client,
) *Server {
	return &Server{
		cfg:   cfg,
		pgDB:  pgDB,
		redis: redis,
		fiber: fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          error_handler.FiberErrorHandler,
		}),
	}
}

func (s *Server) Run() error {
	err := s.MapHandlers()
	if err != nil {
		return err
	}

	go func() {
		s.fiber.Get("/health_check", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
		log.Info("Server is started ", s.cfg.Server.Host)
		err := s.fiber.Listen(s.cfg.Server.Host)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	if err = s.fiber.ShutdownWithContext(context.Background()); err != nil {
		log.Error(err)
	} else {
		log.Info("Fiber closed properly")
	}

	return nil
}

func NewDB(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRedisClient(cfg config.Config) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password,
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrapf(err, "ping")
	}

	return client, nil
}
