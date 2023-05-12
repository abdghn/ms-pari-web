/*
 * Created on 01/04/22 15.31
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package user

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
)

type Usecase interface {
	Create(user *model.User) (*model.User, error)
	ReadAll() (*[]model.User, error)
	ReadById(id int) (*model.User, error)
	Update(id int, user *model.User) (*model.User, error)
	ChangePassword(user request.ChangePassword) (*model.User, error)
	Delete(id int) error
}

type usecase struct {
	repository user.Repository
}

func NewUsecase(repository user.Repository) Usecase {
	return &usecase{repository}
}

func (e *usecase) Create(user *model.User) (*model.User, error) {
	return e.repository.Create(user)
}

func (e *usecase) ReadAll() (*[]model.User, error) {
	return e.repository.ReadAll()
}

func (e *usecase) ReadById(id int) (*model.User, error) {
	return e.repository.ReadById(id)
}

func (e *usecase) Update(id int, user *model.User) (*model.User, error) {
	return e.repository.Update(id, user)
}

func (e *usecase) ChangePassword(changePassword request.ChangePassword) (*model.User, error) {
	userModel, err := e.repository.ReadById(changePassword.UserID)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	userModel.Password = changePassword.Password
	userModel.MustChangePassword = false

	return e.repository.UpdatePasswordLogin(userModel)
}

func (e *usecase) Delete(id int) error {
	return e.repository.Delete(id)
}
