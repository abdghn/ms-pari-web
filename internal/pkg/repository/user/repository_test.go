/*
 * Created on 18/04/22 05.56
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package user

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	gormDb, _ := gorm.Open("mysql", db)
	defer gormDb.Close()

	userRepo := NewRepository(gormDb)

	t.Run("Register", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"name", "password", "email", "role"}).AddRow(
			"Alex", "123456", "alex@gmail.com", 1)

		user := &model.User{
			Name:     "Alex",
			Email:    "alex@gmail.com",
			Password: "123456",
			RoleID:   2,
		}

		mock.ExpectQuery(`INSERT INTO users (name, email, verification_level,status,password, role_id , company_id , created_at,updated_at,deleted_at)
						VALUES ($1, $2, $3, $4, $5, $6, now(), now(), NULL) 
						RETURNING *`).WithArgs(&user.Name, &user.Email,
			&user.Password, &user.RoleID).WillReturnRows(rows)

		createdUser, err := userRepo.Create(user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}
