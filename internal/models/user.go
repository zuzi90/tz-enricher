package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
	"strconv"
	"time"
)

type (
	UserFN struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}

	UserCreate struct {
		Name        string `json:"name,omitempty"        db:"name"`
		Surname     string `json:"surname,omitempty"     db:"surname"`
		Patronymic  string `json:"patronymic,omitempty"  db:"patronymic"`
		Age         int    `json:"age,omitempty"         db:"age"`
		Gender      string `json:"gender,omitempty"      db:"gender"`
		Nationality string `json:"nationality,omitempty" db:"nationality"`
	}

	ResponseFNError struct {
		UserFN
		ErrMessage string `json:"errMessage"`
	}

	User struct {
		ID          int       `db:"id"          json:"id"`
		Name        string    `db:"name"        json:"name"`
		Surname     string    `db:"surname"     json:"surname"`
		Patronymic  string    `db:"patronymic"  json:"patronymic"`
		Age         int       `db:"age"         json:"age"`
		Gender      string    `db:"gender"      json:"gender"`
		Nationality string    `db:"nationality" json:"nationality"`
		IsDeleted   bool      `db:"is_deleted"  json:"isDeleted"`
		CreatedAt   time.Time `db:"created_at"  json:"createdAt"`
		UpdatedAt   time.Time `db:"updated_at"  json:"updatedAt"`
	}

	UserUpdate struct {
		Name        *string `json:"name"        db:"name"`
		Surname     *string `json:"surname"     db:"surname"`
		Patronymic  *string `json:"patronymic"  db:"patronymic"`
		Age         *int    `json:"age"         db:"age"`
		Gender      *string `json:"gender"      db:"gender"`
		Nationality *string `json:"nationality" db:"nationality"`
	}

	GetUsersParams struct {
		Text       string `json:"text"       db:"text"`
		Limit      int    `json:"limit"      db:"limit"`
		Offset     int    `json:"offset"     db:"offset"`
		Sorting    string `json:"sorting"    db:"sorting"`
		Descending bool   `json:"descending" db:"descending"`
	}
)

func (u *UserCreate) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(2, 25)),
		validation.Field(&u.Surname, validation.Required, validation.Length(2, 25)),
		validation.Field(&u.Age, validation.Required, validation.Min(1), validation.Max(120)),
		validation.Field(&u.Gender, validation.Required, validation.Length(2, 25)),
		validation.Field(&u.Nationality, validation.Required, validation.Length(2, 25)))
}

func NewCreateUser(fn UserFN) UserCreate {
	return UserCreate{
		Name:       fn.Name,
		Surname:    fn.Surname,
		Patronymic: fn.Patronymic,
	}
}

func (u *UserFN) ValidateFN() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(2, 25),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")),
		),
		validation.Field(&u.Surname, validation.Required, validation.Length(2, 25),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")),
		),
		validation.Field(&u.Patronymic, validation.Length(2, 25),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$"))),
	)
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

var usersFieldsMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"surname":     "surname",
	"patronymic":  "patronymic",
	"age":         "age",
	"gender":      "gender",
	"nationality": "nationality",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

func NewParams(text string, limit string, offset string, sorting string, descending string) GetUsersParams {
	params := GetUsersParams{}

	params.Text = text
	defaultLimit := 100
	params.Limit, _ = strconv.Atoi(limit)

	if params.Limit == 0 {
		params.Limit = defaultLimit
	}

	params.Offset, _ = strconv.Atoi(offset)
	val, ok := usersFieldsMapping[sorting]
	if !ok {
		params.Sorting = "id"
	} else {
		params.Sorting = val
	}

	if sorting == "" {
		sorting = "id"
	}

	params.Sorting = sorting
	params.Descending, _ = strconv.ParseBool(descending)

	return params
}
