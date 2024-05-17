package userservice

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zuzi90/tz-enricher/internal/models"
)

type userStorage interface {
	CreateUser(ctx context.Context, val models.UserCreate) (*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUsers(ctx context.Context, params models.GetUsersParams) ([]*models.User, error)
	UpdateUser(ctx context.Context, user models.UserUpdate, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type cache interface {
	Get(ctx context.Context, key string) (*models.User, error)
	Set(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, key int) error
}

type UserService struct {
	db    userStorage
	log   *logrus.Entry
	cache cache
}

func NewUserService(db userStorage, logger *logrus.Logger, cache cache) *UserService {
	return &UserService{
		db:    db,
		log:   logger.WithField("module", "user_service"),
		cache: cache,
	}
}

func (s *UserService) CreateUser(ctx context.Context, val models.UserCreate) (*models.User, error) {
	user, err := s.db.CreateUser(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("err creating user db: %w", err)
	}

	if err = s.cache.Set(ctx, user); err != nil {
		return nil, fmt.Errorf("err writing user to cache: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, err := s.cache.Get(ctx, fmt.Sprintf("UserKey:%d", id))
	if err != nil {

		user, err = s.db.GetUser(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("err getting user from db: %w", err)
		}

		err = s.cache.Set(ctx, user)
		switch {
		case err == nil:
			return user, nil
		case err != nil:
			s.log.Warnf("err writing user to cache: %v", err)
			return user, nil
		}
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context, params models.GetUsersParams) ([]*models.User, error) {

	users, err := s.db.GetUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("err failed get users from db: %w", err)
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, val models.UserUpdate) (*models.User, error) {
	user, err := s.db.UpdateUser(ctx, val, id)
	if err != nil {
		return nil, fmt.Errorf("err updating user db: %w", err)
	}

	if err := s.cache.Update(ctx, user); err != nil {
		s.log.Warnf("err updating user in cache: %v", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("err delete user, db delete %w", err)
	}

	err = s.cache.Delete(ctx, id)
	if err != nil {
		s.log.Warnf("err delete user, cache delete: %v", err)
		return nil
	}

	return nil
}
