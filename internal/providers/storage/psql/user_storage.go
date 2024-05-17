package psql

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"github.com/zuzi90/tz-enricher/internal/models"
	"strconv"
)

func (s *Storage) CreateUser(ctx context.Context, val models.UserCreate) (*models.User, error) {
	user := models.User{}

	query := `
			 INSERT INTO users(name, surname, patronymic, age, gender, nationality)
			 VALUES($1,$2,$3,$4,$5,$6)
			 RETURNING id, name, surname, patronymic, age, gender, nationality, is_deleted, created_at, updated_at`
	err := s.db.GetContext(ctx, &user, query, val.Name, val.Surname, val.Patronymic, val.Age, val.Gender, val.Nationality)

	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	query := `
			 SELECT id, name, surname, patronymic, age, gender, nationality,  is_deleted, created_at, updated_at
			 FROM users WHERE id = $1 AND is_deleted = false`

	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return &models.User{}, err
	}

	return &user, nil
}

func (s *Storage) GetUsers(ctx context.Context, params models.GetUsersParams) ([]*models.User, error) {

	var args []interface{}
	var builder bytes.Buffer
	users := make([]*models.User, 0)

	builder.WriteString(`SELECT id, name, surname, patronymic, age, gender, nationality, is_deleted, created_at, updated_at FROM users WHERE is_deleted = false`)

	if params.Text != "" {
		args = append(args, params.Text)
		builder.WriteString(` AND (name LIKE $` + strconv.Itoa(len(args)) + ` OR surname LIKE $` + strconv.Itoa(len(args)) + ` OR patronymic LIKE $` + strconv.Itoa(len(args)) + `)`)
	}

	if params.Sorting != "" {
		builder.WriteString(` ORDER BY ` + params.Sorting)
		if params.Descending {
			builder.WriteString(` DESC`)
		}
	}

	args = append(args, params.Limit)
	builder.WriteString(` LIMIT $` + strconv.Itoa(len(args)))
	args = append(args, params.Offset)
	builder.WriteString(` OFFSET $` + strconv.Itoa(len(args)))
	query := builder.String()

	err := s.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.UserUpdate, id int) (*models.User, error) {
	var args []interface{}

	var builder bytes.Buffer

	var userResponse models.User

	builder.WriteString(`UPDATE users SET updated_at = NOW()`)

	if user.Name != nil {
		args = append(args, *user.Name)
		builder.WriteString(`, name = ` + `$` + strconv.Itoa(len(args)))
	}

	if user.Surname != nil {
		args = append(args, *user.Surname)
		builder.WriteString(`, surname = ` + `$` + strconv.Itoa(len(args)))
	}

	if user.Patronymic != nil {
		args = append(args, *user.Patronymic)
		builder.WriteString(`, patronymic = ` + `$` + strconv.Itoa(len(args)))
	}

	if user.Age != nil {
		args = append(args, *user.Age)
		builder.WriteString(`, age = ` + `$` + strconv.Itoa(len(args)))
	}

	if user.Gender != nil {
		args = append(args, *user.Gender)
		builder.WriteString(`, gender = ` + `$` + strconv.Itoa(len(args)))
	}

	if user.Nationality != nil {
		args = append(args, *user.Nationality)
		builder.WriteString(`, nationality = ` + `$` + strconv.Itoa(len(args)))
	}

	builder.WriteString(`WHERE id =` + strconv.Itoa(id))
	builder.WriteString(`AND is_deleted = false `)
	builder.WriteString(`RETURNING id, name, surname,
							patronymic, age, gender, nationality, is_deleted, created_at, updated_at`)

	err := s.db.QueryRowContext(ctx, builder.String(), args...).Scan(&userResponse.ID,
		&userResponse.Name, &userResponse.Surname,
		&userResponse.Patronymic, &userResponse.Age, &userResponse.Gender, &userResponse.Nationality,
		&userResponse.IsDeleted, &userResponse.CreatedAt, &userResponse.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}

		return &models.User{}, err
	}

	return &userResponse, nil

}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	query := `UPDATE users SET is_deleted = true WHERE id = $1 AND is_deleted = false`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrUserNotFound
	}

	return nil
}
