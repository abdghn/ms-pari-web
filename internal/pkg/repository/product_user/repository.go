/*
 * Created on 27/05/22 00.03
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package product_user

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(productUser *model.ProductUser) (*model.ProductUser, error)
	ReadAll() (*[]model.ProductUser, error)
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.ProductUser, error)
	ReadById(id int) (*model.ProductUser, error)
	ReadBy(criteria map[string]interface{}) (*model.ProductUser, error)
	Update(id int, person *model.ProductUser) (*model.ProductUser, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(product_user *model.ProductUser) (*model.ProductUser, error) {
	err := e.DB.Save(&product_user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return product_user, nil
}

func (e *repository) ReadAll() (*[]model.ProductUser, error) {
	var product_users []model.ProductUser
	err := e.DB.Find(&product_users).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &product_users, nil
}

func (e *repository) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.ProductUser, error) {
	var product_users []model.ProductUser

	query := e.DB.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at DESC").Limit(limit).Find(&product_users).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &product_users, nil
}

func (e *repository) ReadById(id int) (*model.ProductUser, error) {
	var product_user = model.ProductUser{}
	err := e.DB.Table("product_users").Where("id = ?", id).First(&product_user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &product_user, nil
}

func (e *repository) Update(id int, product_user *model.ProductUser) (*model.ProductUser, error) {
	var upProductUser = model.ProductUser{}
	err := e.DB.Table("product_users").Where("id = ?", id).First(&upProductUser).Update(&product_user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upProductUser, nil
}

func (e *repository) Delete(id int) error {
	var product_user = model.ProductUser{}
	err := e.DB.Table("product_users").Where("id = ?", id).First(&product_user).Delete(&product_user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *repository) Count(criteria map[string]interface{}) int {
	var result int
	err := e.DB.Table("product_users").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (e *repository) ReadBy(criteria map[string]interface{}) (*model.ProductUser, error) {
	var product_user = model.ProductUser{}
	err := e.DB.Table("product_users").Where(criteria).First(&product_user).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productUserRepository.ReadBy] error execute query %v \n", err)
		return nil, err
	}
	return &product_user, nil
}
