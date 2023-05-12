/*
 * Created on 01/04/22 15.32
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package user

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(user *model.User) (*model.User, error)
	ReadAll() (*[]model.User, error)
	ReadById(id int) (*model.User, error)
	ReadByEmail(email string) (*model.User, error)
	Update(id int, user *model.User) (*model.User, error)
	UpdatePasswordLogin(user *model.User) (*model.User, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(user *model.User) (*model.User, error) {
	//var role = model.Role{}
	//var company = model.Company{}

	tx := e.DB.Begin()
	defer tx.Rollback()

	err := tx.Save(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}

	//err = tx.Table("roles").Where("id = ?", user.RoleID).First(&role).Error
	//if err != nil {
	helper.CommonLogger().Error(err)
	//	fmt.Printf("[repository.Create] error execute query %v \n", err)
	//	return nil, fmt.Errorf("id role is not exists")
	//}
	//
	//err = tx.Table("companies").Where("id = ?", user.CompanyID).First(&company).Error
	//if err != nil {
	helper.CommonLogger().Error(err)
	//	fmt.Printf("[repository.Create] error execute query %v \n", err)
	//	return nil, fmt.Errorf("id company is not exists")
	//}
	//
	//user.Role = role
	//user.Company = company

	tx.Commit()

	return user, nil
}

func (e *repository) ReadAll() (*[]model.User, error) {
	var users []model.User
	err := e.DB.Find(&users).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &users, nil
}

func (e *repository) ReadById(id int) (*model.User, error) {
	var user = model.User{}
	err := e.DB.Select("*, r.name AS role_name, c.name AS company_name").
		Table("users").
		Joins("JOIN roles r ON r.id = users.role_id").
		Joins("JOIN companies c ON c.id = users.company_id").
		Where("users.id = ?", id).First(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &user, nil
}

func (e *repository) ReadByEmail(email string) (*model.User, error) {
	var user = model.User{}
	err := e.DB.Select("*, r.name AS role_name, c.name AS company_name").
		Table("users").
		Joins("JOIN roles r ON r.id = users.role_id").
		Joins("JOIN companies c ON c.id = users.company_id").
		Where("users.email = ?", email).First(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadByEmail] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &user, nil
}

func (e *repository) Update(id int, user *model.User) (*model.User, error) {
	var upUser = model.User{}
	err := e.DB.Table("users").Where("id = ?", id).First(&upUser).Update(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upUser, nil
}

func (e *repository) UpdatePasswordLogin(user *model.User) (*model.User, error) {
	var upUser = model.User{}
	err := e.DB.Model(&user).Update(map[string]interface{}{"password": &user.Password, "must_change_password": &user.MustChangePassword}).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upUser, nil
}

func (e *repository) Delete(id int) error {
	var user = model.User{}
	err := e.DB.Table("users").Where("id = ?", id).First(&user).Delete(&user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *repository) Count(criteria map[string]interface{}) int {
	var result int
	err := e.DB.Table("users").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}
