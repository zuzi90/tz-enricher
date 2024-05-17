package rest

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/zuzi90/tz-enricher/internal/models"
	"net/http"
)

type messageService interface {
	Handle(ctx context.Context, val []byte) error
}

type userService interface {
	CreateUser(ctx context.Context, val models.UserCreate) (*models.User, error)
	DeleteUser(ctx context.Context, key int) error
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUsers(ctx context.Context, params models.GetUsersParams) ([]*models.User, error)
	UpdateUser(ctx context.Context, id int, val models.UserUpdate) (*models.User, error)
}

type Server struct {
	log      *logrus.Entry
	router   *chi.Mux
	server   *http.Server
	address  string
	services messageService
	uService userService
}

func NewServer(port string, log *logrus.Logger, services messageService, uService userService) *Server {
	srv := Server{
		log:      log.WithField("module", "server"),
		router:   chi.NewRouter(),
		address:  port,
		services: services,
		uService: uService,
	}

	srv.InitRoutes()

	server := http.Server{
		Addr:    srv.address,
		Handler: srv.router,
	}

	srv.server = &server

	return &srv
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		if err := s.server.Shutdown(ctx); err != nil {
			s.log.Info("Closing HTTP server", err)
		}
	}()

	s.log.Infof("Listening on %v", s.address)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) response(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.log.Warnf("err encoding dat %v:", err)
			return
		}
	}
}

func (s *Server) responseOk(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
