package tests

import (
	"context"
	"encoding/json"
	"github.com/zuzi90/tz-enricher/internal/models"
	"net/http"
	"strconv"
)

func (s *IntegrationTestSuite) TestCreateUser() {

	ctx := context.Background()

	val := models.UserCreate{
		Name:        "Golda",
		Surname:     "Meir",
		Patronymic:  "Meyerson",
		Age:         80,
		Gender:      "female",
		Nationality: "Israelite",
	}

	reqBody, err := json.Marshal(val)
	s.Require().NoError(err)

	s.Run("create user", func() {
		var userResp models.User
		s.userID = userResp.ID
		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusCreated, code)
	})

	val2 := models.UserCreate{
		Name:        "Golda",
		Surname:     "",
		Patronymic:  "Meyerson",
		Age:         80,
		Gender:      "female",
		Nationality: "Israelite",
	}

	reqBody, err = json.Marshal(val2)
	s.Require().NoError(err)

	s.Run("create user with empty fields", func() {
		var userResp models.User

		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusBadRequest, code)
	})
}

func (s *IntegrationTestSuite) TestGetUser() {
	ctx := context.Background()

	val := models.UserCreate{
		Name:        "Golda",
		Surname:     "Meir",
		Patronymic:  "Meyerson",
		Age:         80,
		Gender:      "female",
		Nationality: "Israelite",
	}

	reqBody, err := json.Marshal(val)
	s.Require().NoError(err)

	s.Run("create user", func() {
		var userResp models.User
		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusCreated, code)
		s.userID = userResp.ID
	})

	s.Run("get user", func() {
		var userResp models.User
		code := s.sendRequest(s.T(), ctx, http.MethodGet, s.host, "/api/v1/users/"+strconv.Itoa(s.userID), []byte{}, &userResp, nil)
		s.Require().Equal(http.StatusOK, code)
		s.Require().Equal("Golda", userResp.Name)
	})

	s.Run("get a invalid id", func() {
		var userResp models.User
		invalidID := "A99T"
		code := s.sendRequest(s.T(), ctx, http.MethodGet, s.host, "/api/v1/users/"+invalidID, []byte{}, &userResp, nil)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	s.Run("get a non-existent id", func() {
		var userResp models.User
		nonExistID := "99999"
		code := s.sendRequest(s.T(), ctx, http.MethodGet, s.host, "/api/v1/users/"+nonExistID, []byte{}, &userResp, nil)
		s.Require().Equal(http.StatusNotFound, code)
	})
}

func (s *IntegrationTestSuite) TestGetUsers() {
	ctx := context.Background()

	val := models.UserCreate{
		Name:        "Shoshana",
		Surname:     "Valerievna",
		Patronymic:  "Ivanov",
		Age:         40,
		Gender:      "female",
		Nationality: "Israelite",
	}

	reqBody, err := json.Marshal(val)
	s.Require().NoError(err)

	s.Run("create a user with name Shoshana", func() {
		var userResp models.User
		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusCreated, code)
	})

	text := "Shoshana"
	limit := "10"
	offset := "0"
	sorting := "age"
	descending := "false"

	params := models.NewParams(text, limit, offset, sorting, descending)
	var usersResp []*models.User
	s.Run("get users with name Shoshana", func() {
		code := s.sendRequest(s.T(), ctx, http.MethodGet, s.host, "/api/v1/users/", []byte{}, &usersResp, &params)
		s.Require().Equal(http.StatusOK, code)
		s.Require().Equal("Shoshana", usersResp[0].Name)
	})

}

func (s *IntegrationTestSuite) TestUpdateUser() {
	ctx := context.Background()

	val := models.UserCreate{
		Name:        "Moshe",
		Surname:     "Sharett",
		Patronymic:  "Jacob",
		Age:         71,
		Gender:      "male",
		Nationality: "Israeli",
	}

	reqBody, err := json.Marshal(val)
	s.Require().NoError(err)

	s.Run("create user", func() {
		var userResp models.User
		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusCreated, code)
		s.userID = userResp.ID

	})

	s.Run("update user", func() {
		val1 := models.UserCreate{
			Name:        "Irina",
			Nationality: "Israelite",
		}

		reqBody, err = json.Marshal(val1)
		s.Require().NoError(err)

		var userResp models.User

		code := s.sendRequest(s.T(), ctx, http.MethodPatch, s.host, "/api/v1/users/"+strconv.Itoa(s.userID), reqBody, &userResp, nil)
		s.Require().Equal(http.StatusOK, code)
		s.Require().Equal("Irina", userResp.Name)
	})

	s.Run("update user with non exist id", func() {
		val1 := models.UserCreate{
			Nationality: "Israelite",
		}

		reqBody, err = json.Marshal(val1)
		s.Require().NoError(err)

		var userResp models.User

		nonExistID := "99999"
		code := s.sendRequest(s.T(), ctx, http.MethodPatch, s.host, "/api/v1/users/"+nonExistID, reqBody, &userResp, nil)
		s.Require().Equal(http.StatusNotFound, code)
	})

	s.Run("update user with invalid id", func() {
		val1 := models.UserCreate{
			Nationality: "Dagestanit",
		}

		reqBody, err = json.Marshal(val1)
		s.Require().NoError(err)

		var userResp models.User

		invalidID := "A99T"

		code := s.sendRequest(s.T(), ctx, http.MethodPatch, s.host, "/api/v1/users/"+invalidID, reqBody, &userResp, nil)
		s.Require().Equal(http.StatusBadRequest, code)
	})
}

func (s *IntegrationTestSuite) TestDeleteUser() {
	ctx := context.Background()

	val := models.UserCreate{
		Name:        "David",
		Surname:     "Ben-Gurion",
		Patronymic:  "Avigdor",
		Age:         87,
		Gender:      "male",
		Nationality: "israeli",
	}

	reqBody, err := json.Marshal(val)
	s.Require().NoError(err)

	s.Run("create user", func() {
		var userResp models.User
		code := s.sendRequest(s.T(), ctx, http.MethodPost, s.host, "/api/v1/users", reqBody, &userResp, nil)
		s.Require().Equal(http.StatusCreated, code)
		s.userID = userResp.ID
	})

	s.Run("delete user", func() {
		code := s.sendRequest(s.T(), ctx, http.MethodDelete, s.host, "/api/v1/users/"+strconv.Itoa(s.userID), []byte{}, nil, nil)
		s.Require().Equal(http.StatusOK, code)
	})

	s.Run("delete a invalid id", func() {
		invalidID := "A99T"
		code := s.sendRequest(s.T(), ctx, http.MethodDelete, s.host, "/api/v1/users/"+invalidID, []byte{}, nil, nil)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	s.Run("delete a non-existent id", func() {
		nonExistID := "99999"
		code := s.sendRequest(s.T(), ctx, http.MethodDelete, s.host, "/api/v1/users/"+nonExistID, []byte{}, nil, nil)
		s.Require().Equal(http.StatusNotFound, code)
	})

}
