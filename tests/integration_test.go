package tests

import (
	"bytes"
	"context"
	"encoding/json"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/zuzi90/tz-enricher/internal/config"
	"github.com/zuzi90/tz-enricher/internal/logger"
	"github.com/zuzi90/tz-enricher/internal/models"
	"github.com/zuzi90/tz-enricher/internal/providers/cache"
	"github.com/zuzi90/tz-enricher/internal/providers/kafka"
	"github.com/zuzi90/tz-enricher/internal/providers/storage/psql"
	"github.com/zuzi90/tz-enricher/internal/rest"
	message_service "github.com/zuzi90/tz-enricher/internal/services/message-service"
	"github.com/zuzi90/tz-enricher/internal/services/resolvers"
	"github.com/zuzi90/tz-enricher/internal/services/userservice"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const (
	pgDSN    = "postgresql://postgres:secret@localhost:5431/tzDB?sslmode=disable"
	port     = ":5005"
	logLevel = "debug"
)

type IntegrationTestSuite struct {
	conf            *config.Config
	client          *http.Client
	log             *logrus.Logger
	db              *psql.Storage
	service         *message_service.MessageService
	uService        *userservice.UserService
	server          *rest.Server
	cancel          context.CancelFunc
	cache           *cache.Redis
	consumer        *kafka.Consumer
	producer        *kafka.Producer
	ageResolver     *resolvers.AgeResolver
	genderResolver  *resolvers.GenderResolver
	countryResolver *resolvers.CountryResolver
	userID          int
	pgDSN           string
	host            string

	suite.Suite
}

func (s *IntegrationTestSuite) SetupSuite() {
	var err error
	s.userID = 0
	s.pgDSN = pgDSN
	s.host = "http://localhost" + port

	s.conf, err = config.NewConfig()
	s.Require().NoError(err)

	s.log, err = logger.NewLogger(logLevel)
	s.Require().NoError(err)

	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())

	s.client = &http.Client{}

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     s.conf.RedisDSN,
		Password: s.conf.RedisPassword,
		DB:       s.conf.RedisDB,
	})

	status := clientRedis.Ping(ctx)
	s.Require().NoError(status.Err())
	s.log.Warnf("redis ping: %v", err)

	s.cache = cache.NewRedis(clientRedis, s.log)

	s.db, err = psql.NewStorage(s.log, s.pgDSN)
	s.Require().NoError(err)
	err = s.db.UpdateSchema()
	s.Require().NoError(err)

	s.producer, err = kafka.NewProducer(s.conf.Brokers, s.conf.KafkaTopicWrongFN, s.log)
	s.Require().NoError(err)

	s.ageResolver = resolvers.NewAgeResolver(s.log, s.conf.AgeURL)
	s.genderResolver = resolvers.NewGenderResolver(s.log, s.conf.GenderURL)
	s.countryResolver = resolvers.NewCountryResolver(s.log, s.conf.NationalityURL)

	s.service = message_service.NewMessageService(s.cache, s.ageResolver, s.genderResolver, s.countryResolver, s.producer, s.db)
	s.uService = userservice.NewUserService(s.db, s.log, s.cache)

	s.consumer = kafka.NewConsumer(s.conf.Brokers, s.conf.WorkersCount, s.conf.KafkaTopic, s.log, s.service)

	s.server = rest.NewServer(port, s.log, s.service, s.uService)

	go func() {
		err = s.consumer.Run(ctx)
		s.Require().NoError(err)
	}()

	go func() {
		err = s.server.Run(ctx)
		s.Require().NoError(err)
	}()

	time.Sleep(800 * time.Millisecond)
}

func (s *IntegrationTestSuite) sendRequest(t *testing.T, ctx context.Context, method, host, endpoint string, requestBody []byte, result any, params *models.GetUsersParams) int {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, method, host+endpoint, bytes.NewReader(requestBody))

	s.Require().NoError(err)

	if params != nil {
		query := req.URL.Query()
		query.Set("text", params.Text)
		query.Set("limit", strconv.Itoa(params.Limit))
		query.Set("offset", strconv.Itoa(params.Offset))
		query.Set("sorting", params.Sorting)
		query.Set("descending", strconv.FormatBool(params.Descending))
		req.URL.RawQuery = query.Encode()
	}

	s.Require().NoError(err)

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/json" {
		return resp.StatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	s.Require().NoError(err)

	return resp.StatusCode
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.cancel()
}

func (s *IntegrationTestSuite) TearDownTest() {
	err := s.db.TruncateTables(`users`)
	s.Require().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
