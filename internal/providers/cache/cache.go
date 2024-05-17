package cache

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/zuzi90/tz-enricher/internal/models"
)

var ErrWritingCache = errors.New("error writing Redis cache")
var ErrUpdatingCache = errors.New("error updating Redis cache")
var ErrDeletingCache = errors.New("error deleting from Redis cache")
var ErrUserNotFound = errors.New("error user not found in Redis cache")

type Redis struct {
	client *redis.Client
	log    *logrus.Entry
}

func NewRedis(client *redis.Client, log *logrus.Logger) *Redis {
	return &Redis{
		client: client,
		log:    log.WithField("module", "server"),
	}
}

func (r *Redis) Get(ctx context.Context, key string) (*models.User, error) {

	val, err := r.client.Get(ctx, "UserKey:"+key).Result()
	if err == redis.Nil {
		return &models.User{}, ErrUserNotFound
	}

	if err != nil {
		return &models.User{}, err
	}

	var user *models.User

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal([]byte(val), &user); err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (r *Redis) Set(ctx context.Context, user *models.User) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	userByte, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, fmt.Sprintf("UserKey:%d", user.ID), userByte, 0).Err()
	if err != nil {
		return ErrWritingCache
	}

	return nil
}

func (r *Redis) Update(ctx context.Context, user *models.User) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	userByte, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, fmt.Sprintf("UserKey:%d", user.ID), userByte, 0).Err()
	if err != nil {
		return ErrUpdatingCache
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, key int) error {
	_, err := r.client.Del(ctx, fmt.Sprintf("UserKey:%d", key)).Result()

	if err != nil {
		return ErrDeletingCache
	}

	return nil
}
