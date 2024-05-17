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

type AgeResolver struct {
	client *http.Client
	log    *logrus.Entry
	ageURL string
}

func NewAgeResolver(log *logrus.Logger, url string) *AgeResolver {
	resolver := AgeResolver{
		log:    log.WithField("module", "AResolver"),
		client: &http.Client{Timeout: 5 * time.Second},
		ageURL: url,
	}

	return &resolver
}

func (r *AgeResolver) GetAge(ctx context.Context, name string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.ageURL+name, nil)
	if err != nil {
		return 0, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("age resolver: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			r.log.Warnf("get user age, closing response body err: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		r.log.Warnf("unexpected status code %d", resp.StatusCode)
		return 0, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	age := models.AgeResolver{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &age); err != nil {
		return 0, err
	}

	if age.Age < 0 || age.Age == 0 {
		return 0, fmt.Errorf("age is negative or equal to zero")
	}

	return age.Age, nil
}
