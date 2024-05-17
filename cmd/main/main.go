package main

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/zuzi90/tz-enricher/internal/config"
	"github.com/zuzi90/tz-enricher/internal/logger"
	"github.com/zuzi90/tz-enricher/internal/providers/cache"
	"github.com/zuzi90/tz-enricher/internal/providers/kafka"
	"github.com/zuzi90/tz-enricher/internal/providers/storage/psql"
	"github.com/zuzi90/tz-enricher/internal/rest"
	message_service "github.com/zuzi90/tz-enricher/internal/services/message-service"
	"github.com/zuzi90/tz-enricher/internal/services/resolvers"
	"github.com/zuzi90/tz-enricher/internal/services/userservice"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
)

// @title Resolver API
// @version 1.0
// @description API server Resolver

// @host localhost:5005
// @BasePath /

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	log, err := logger.NewLogger(cfg.LogLvl)
	if err != nil {
		return err
	}

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM)
	defer cancel()

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDSN,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	status := clientRedis.Ping(ctx)

	if status.Err() != nil {
		log.Warnf("redis ping: %v", err)
		return status.Err()
	}

	log.Info("Redis up")
	rCache := cache.NewRedis(clientRedis, log)

	db, err := psql.NewStorage(log, cfg.PgDSN)
	if err != nil {
		return err
	}

	defer func() {
		if err = db.CloseDB(); err != nil {
			log.Warnf("closing db: %v", err)
		}
	}()

	if err = db.UpdateSchema(); err != nil {
		log.Warnf("update schema: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Brokers, cfg.KafkaTopicWrongFN, log)
	if err != nil {
		return err
	}

	ageResolver := resolvers.NewAgeResolver(log, cfg.AgeURL)
	genderResolver := resolvers.NewGenderResolver(log, cfg.GenderURL)
	countryResolver := resolvers.NewCountryResolver(log, cfg.NationalityURL)

	mService := message_service.NewMessageService(rCache, ageResolver, genderResolver, countryResolver, producer, db)
	uService := userservice.NewUserService(db, log, rCache)

	consumer := kafka.NewConsumer(cfg.Brokers, cfg.WorkersCount, cfg.KafkaTopic, log, mService)

	server := rest.NewServer(cfg.ServerPORT, log, mService, uService)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return consumer.Run(ctx)
	})

	eg.Go(func() error {
		return server.Run(ctx)
	})

	if err = eg.Wait(); err != nil {
		return err
	}

	return nil
}
