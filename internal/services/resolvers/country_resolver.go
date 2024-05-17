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

type CountryResolver struct {
	log        *logrus.Entry
	client     *http.Client
	countryURL string
}

func NewCountryResolver(log *logrus.Logger, url string) *CountryResolver {
	resolver := CountryResolver{
		log:        log.WithField("module", "CountryResolver"),
		client:     &http.Client{Timeout: 5 * time.Second},
		countryURL: url,
	}

	return &resolver
}

func (r *CountryResolver) GetCountry(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.countryURL+name, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("country resolver: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			r.log.Warnf("get user country, closing response body err: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		r.log.Warnf("unexpected status code %d", resp.StatusCode)
		return "", fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	country := models.NationalityResolver{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &country); err != nil {
		return "", err
	}

	if len(country.Country) == 0 || country.Country[0].CountryID == "" {
		return "", fmt.Errorf("country is empty, name: %s", name)
	}

	return country.Country[0].CountryID, nil
}
