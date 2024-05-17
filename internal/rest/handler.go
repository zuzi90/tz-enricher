package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/zuzi90/tz-enricher/internal/models"
	"net/http"
	"strconv"
)

// @Summary Создать пользователя
// @Tags user
// @Description create user
// @Accept json
// @Produce json
// @Param input body models.UserCreate true "account info"
// @Success 201 {object} models.User
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users [post].
func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userReq := models.UserCreate{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		s.log.Warnf("err encoding dat %v:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := userReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.uService.CreateUser(ctx, userReq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.log.WithField("request", userReq).Warnf("err creating user: %v", err)
		return
	}

	s.response(w, http.StatusCreated, user)
}

// @Summary Получить пользователя по id
// @Tags user
// @Description get user
// @Produce json
// @Param id  path  string  true  "id"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users/{id} [get].
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := chi.URLParam(r, "id")

	id, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.uService.GetUser(ctx, id)
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.log.WithField("request", id).Warnf("err getting user: %v", err)
		return
	}

	s.response(w, http.StatusOK, user)
}

// @Summary Получить список пользователей
// @Tags user
// @Description get users
// @Produce json
// @Param text query string false "text"
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param sorting query string false "sorting"
// @Param descending query string false "descending"
// @Success 200 {array} models.User
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users/ [get].
func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	text := r.URL.Query().Get("text")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	sorting := r.URL.Query().Get("sorting")
	descending := r.URL.Query().Get("descending")

	params := models.NewParams(text, limit, offset, sorting, descending)

	users, err := s.uService.GetUsers(ctx, params)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.log.WithField("request", "getUsers").Warnf("err get all users: %v", err)
		return

	}

	s.response(w, http.StatusOK, users)
}

// @Summary Обновить пользователя
// @Tags user
// @Description update user
// @Accept json
// @Produce json
// @Param id  path  string  true  "id"
// @Param input body models.UserUpdate true "account info"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users/{id} [patch].
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := chi.URLParam(r, "id")
	id, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userReq := models.UserUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.uService.UpdateUser(ctx, id, userReq)
	switch {
	case errors.Is(err, models.ErrNoRows):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.log.WithField("request", id).Warnf("err updating user: %v", err)
		return
	}

	s.response(w, http.StatusOK, user)
}

// @Summary Удалить пользователя
// @Tags user
// @Description delete user
// @Produce json
// @Param id  path  string  true  "id"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/users/{id} [delete].
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := chi.URLParam(r, "id")

	id, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.uService.DeleteUser(ctx, id)

	switch {
	case errors.Is(err, models.ErrUserNotFound):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.log.WithField("request", id).Warnf("err deleting user: %v", err)
		return
	}

	s.responseOk(w, http.StatusOK)
}
