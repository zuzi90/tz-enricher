package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := UserCreate{
			Name:        "Frodo",
			Surname:     "Baggins",
			Patronymic:  "Gendolfovich",
			Age:         18,
			Gender:      "male",
			Nationality: "hobbit",
		}
		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("success, empty Patronymic", func(t *testing.T) {
		user := UserCreate{
			Name:        "Frodo",
			Surname:     "Baggins",
			Patronymic:  "",
			Age:         18,
			Gender:      "male",
			Nationality: "hobbit",
		}
		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("failure, empty Surname", func(t *testing.T) {
		user := UserCreate{
			Name:        "Frodo",
			Surname:     "",
			Patronymic:  "Gendolfovich",
			Age:         18,
			Gender:      "male",
			Nationality: "hobbit",
		}
		err := user.Validate()
		assert.Error(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		user := UserCreate{}
		err := user.Validate()
		assert.Error(t, err)
	})

}

func Test_ValidateFIO(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := UserFN{
			Name:       "Frodo",
			Surname:    "Baggins",
			Patronymic: "Gendolfovich",
		}
		err := user.ValidateFN()
		assert.NoError(t, err)
	})

	t.Run("success", func(t *testing.T) {
		user := UserFN{
			Name:       "Frodo",
			Surname:    "Baggins",
			Patronymic: "",
		}
		err := user.ValidateFN()
		assert.NoError(t, err)
	})

	t.Run("failure, empty Surname", func(t *testing.T) {
		user := UserFN{
			Name:       "Frodo",
			Surname:    "",
			Patronymic: "Gendolfovich",
		}
		err := user.ValidateFN()
		assert.Error(t, err)
	})

	t.Run("failure, empty Name", func(t *testing.T) {
		user := UserFN{
			Name:       "",
			Surname:    "Baggins",
			Patronymic: "Gendolfovich",
		}
		err := user.ValidateFN()
		assert.Error(t, err)
	})

	t.Run("failure, forbidden symbol", func(t *testing.T) {
		user := UserFN{
			Name:       "Frodo1",
			Surname:    "Baggins",
			Patronymic: "Gendolfovich",
		}
		err := user.ValidateFN()
		assert.Error(t, err)
	})
}
