package banners_repository

import (
	"avito/assignment/config"
	"avito/assignment/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"strconv"
	"strings"
	"time"
)

type ClientRedisRepo struct {
	db  *redis.Client
	cfg *config.Config
}

func NewClientRedisRepository(db *redis.Client, cfg *config.Config) *ClientRedisRepo {
	return &ClientRedisRepo{
		db:  db,
		cfg: cfg,
	}
}

func (r *ClientRedisRepo) PutBannerRedis(ctx context.Context, putRedisBannerParams *PutRedisBanner) error {
	ctx, span := otel.Tracer("").Start(ctx, "ClientRedisRepo.PutBanner")
	defer span.End()

	sessionBytes, err := json.Marshal(putRedisBannerParams)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ClientRedisRepo.PutBannerRedis.Marshal; err = %s", err.Error()))
	}

	for _, tagId := range putRedisBannerParams.TagIds {
		key := r.createDbKey(tagId, putRedisBannerParams.FeatureId)
		_, err = r.db.Set(ctx, key, sessionBytes, time.Duration(r.cfg.BannerSettings.BannerTTLSeconds)*time.Second).Result()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ClientRedisRepo.PutBannerRedis.Set; err = %s", err.Error()))
		}
	}

	return nil
}

func (r *ClientRedisRepo) GetBannerRedis(ctx context.Context, featureId models.FeatureId, tagId models.TagId) (*models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ClientRedisRepo.GetBannerRedis")
	defer span.End()

	result := &models.FullBanner{}
	key := r.createDbKey(tagId, featureId)

	valueString, err := r.db.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, fiber.ErrNotFound
	} else if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ClientRedisRepo.GetBannerRedis.Get; err = %s", err.Error()))
	}

	err = json.Unmarshal([]byte(valueString), &result)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ClientRedisRepo.GetBannerRedis.Unmarshal; err = %s", err.Error()))
	}

	return result, nil
}

func (r *ClientRedisRepo) DelBannerRedis(ctx context.Context, featureId models.FeatureId, tagId models.TagId) error {
	ctx, span := otel.Tracer("").Start(ctx, "AdminRedisRepo.DelSession")
	defer span.End()

	regex := r.createDbKey(tagId, featureId)
	keys := make([]string, 0, 10)

	iter := r.db.Scan(ctx, 0, regex, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if len(keys) == 0 {
		return nil
	}

	_, err := r.db.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRedisRepo) createDbKey(tagId models.TagId, featureId models.FeatureId) string {
	return strings.Join([]string{strconv.Itoa(int(tagId)), strconv.Itoa(int(featureId))}, ":")
}
