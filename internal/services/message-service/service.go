package message_service

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/zuzi90/tz-enricher/internal/models"
	"golang.org/x/sync/errgroup"
	"time"
)

type messageProducer interface {
	SendMessage(msg []byte) error
}

type ageResolver interface {
	GetAge(ctx context.Context, name string) (int, error)
}

type genderResolver interface {
	GetGender(ctx context.Context, name string) (string, error)
}

type countryResolver interface {
	GetCountry(ctx context.Context, name string) (string, error)
}

type appStorage interface {
	CreateUser(ctx context.Context, user models.UserCreate) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type cache interface {
	Set(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, key int) error
	Update(ctx context.Context, user *models.User) error
}

type MessageService struct {
	metrics         *metrics
	cache           cache
	ageResolver     ageResolver
	genderResolver  genderResolver
	countryResolver countryResolver
	messageProducer messageProducer
	db              appStorage
}

func NewMessageService(
	cache cache,
	ageResolver ageResolver,
	genderResolver genderResolver,
	countryResolver countryResolver,
	messageProducer messageProducer,
	db appStorage,
) *MessageService {
	return &MessageService{
		metrics:         newMetrics(),
		cache:           cache,
		ageResolver:     ageResolver,
		genderResolver:  genderResolver,
		countryResolver: countryResolver,
		messageProducer: messageProducer,
		db:              db,
	}
}

func (s *MessageService) Handle(ctx context.Context, val []byte) error {

	started := time.Now()
	defer func() {
		s.metrics.observe(time.Since(started))
	}()

	fn := models.UserFN{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(val, &fn); err != nil {
		return err
	}

	if err := fn.ValidateFN(); err != nil {
		s.metrics.incInvalidFN(err)
		resp := models.ResponseFNError{}
		resp.UserFN = fn
		resp.ErrMessage = err.Error()

		json := jsoniter.ConfigCompatibleWithStandardLibrary
		respByte, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		if err := s.messageProducer.SendMessage(respByte); err != nil {
			return err
		}

		return nil
	}

	result := models.NewCreateUser(fn)
	eg, ctxE := errgroup.WithContext(ctx)
	eg.Go(func() error {
		age, err := s.ageResolver.GetAge(ctxE, fn.Name)
		if err != nil {
			return err
		}

		result.Age = age

		return nil
	})

	eg.Go(func() error {
		gender, err := s.genderResolver.GetGender(ctxE, fn.Name)
		if err != nil {
			return err
		}

		result.Gender = gender

		return nil
	})

	eg.Go(func() error {
		country, err := s.countryResolver.GetCountry(ctxE, fn.Name)
		if err != nil {
			return err
		}

		result.Nationality = country

		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("eg.Wait(): %w", err)
	}

	user, err := s.db.CreateUser(ctx, result)
	if err != nil {
		return fmt.Errorf("err db create user %w", err)
	}

	if err := s.cache.Set(ctx, user); err != nil {
		return fmt.Errorf("err cache set %w", err)
	}

	return nil
}

func (s *MessageService) CreateUser(ctx context.Context, val models.UserCreate) (*models.User, error) {
	user, err := s.db.CreateUser(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("err db create user %w", err)
	}

	if err := s.cache.Set(ctx, user); err != nil {
		return user, fmt.Errorf("err create user, cache set %w", err)
	}

	return user, nil
}

func (s *MessageService) DeleteUser(ctx context.Context, id int) error {
	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("err db delete user %w", err)
	}

	err = s.cache.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("err get user, cache get %w", err)
	}

	return nil
}
