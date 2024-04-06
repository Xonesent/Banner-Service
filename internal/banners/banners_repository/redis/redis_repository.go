package banners_redis

import (
	"avito/assignment/config"
	"avito/assignment/internal/models/banner_models"
	"avito/assignment/pkg/traces"
	"context"
	"encoding/json"
	"errors"
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

func (r *ClientRedisRepo) PutBanner(ctx context.Context, params banner_models.PutRedisBanner) error {
	ctx, span := otel.Tracer("").Start(ctx, "ClientRedisRepo.PutBanner")
	defer span.End()

	sessionBytes, err := json.Marshal(params.Content)
	if err != nil {
		return traces.SpanSetErrWrapf(
			span,
			err,
			"ClientRedisRepo.PutBanner.Marshal(args: %v)",
			params,
		)
	}

	for _, tagId := range params.TagIds {
		key := r.createDbKey(tagId, params.FeatureId)
		_, err = r.db.Set(ctx, key, sessionBytes, time.Duration(r.cfg.BannerSettings.BannerTTLSeconds)*time.Second).Result()
		if err != nil {
			return traces.SpanSetErrWrapf(
				span,
				err,
				"ClientRedisRepo.PutBanner.Set(args: %v)",
				key,
			)
		}
	}

	return nil
}

func (r *ClientRedisRepo) GetBanner(ctx context.Context, params banner_models.GetRedisBanner) (*banner_models.BannerContent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ClientRedisRepo.GetBanner")
	defer span.End()

	result := &banner_models.BannerContent{}
	key := r.createDbKey(params.TagId, params.FeatureId)

	valueString, err := r.db.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, traces.SpanSetErrWrapf(
			span,
			fiber.ErrNotFound,
			"ClientRedisRepo.GetBanner.redis.Nil(args: %v)",
			key,
		)
	} else if err != nil {
		return nil, traces.SpanSetErrWrapf(
			span,
			err,
			"ClientRedisRepo.GetBanner.Get(args: %v)",
			key,
		)
	}

	err = json.Unmarshal([]byte(valueString), &result)
	if err != nil {
		return nil, traces.SpanSetErrWrapf(
			span,
			err,
			"ClientRedisRepo.GetBanner.Unmarshal(args: %v)",
			valueString,
		)
	}

	return result, nil
}

func (r *ClientRedisRepo) createDbKey(tagId, featureId int) string {
	return strings.Join([]string{strconv.Itoa(tagId), strconv.Itoa(featureId)}, ":")
}
