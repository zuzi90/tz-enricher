package resolvers

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/zuzi90/tz-enricher/internal/models"
	"io"
	"net/http"
	"time"
)

type GenderResolver struct {
	client    *http.Client
	log       *logrus.Entry
	genderURL string
}

func NewGenderResolver(log *logrus.Logger, url string) *GenderResolver {
	resolver := GenderResolver{
		log:       log.WithField("module", "GenderResolver"),
		client:    &http.Client{Timeout: 5 * time.Second},
		genderURL: url,
	}

	return &resolver
}

func (r *GenderResolver) GetGender(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.genderURL+name, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf(" err gender resolver: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			r.log.Warnf("get user gender, closing response body err: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		r.log.Warnf("unexpected status code %d", resp.StatusCode)
		return "", fmt.Errorf("err response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	gender := models.GenderResolver{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &gender); err != nil {
		return "", err
	}

	if gender.Gender == "" {
		return "", fmt.Errorf("err gender is empty, name: %s", name)
	}

	return gender.Gender, nil
}
