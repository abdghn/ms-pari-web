package giro

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(person *model.Giro) (*model.Giro, error)
	ReadAll() (*[]model.Giro, error)
	ReadById(id int) (*model.Giro, error)
	ReadByCode(code string) (*model.Giro, error)
	Update(id int, person *model.Giro) (*model.Giro, error)
	Delete(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(giro *model.Giro) (*model.Giro, error) {
	err := e.DB.Save(&giro).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return giro, nil
}

func (e *repository) ReadAll() (*[]model.Giro, error) {
	var giros []model.Giro
	err := e.DB.Find(&giros).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &giros, nil
}

func (e *repository) ReadById(id int) (*model.Giro, error) {
	var giro = model.Giro{}
	err := e.DB.Table("giros").Where("id = ?", id).First(&giro).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &giro, nil
}
func (e *repository) ReadByCode(code string) (*model.Giro, error) {
	var giro = model.Giro{}
	err := e.DB.Table("giros").Where("code = ?", code).First(&giro).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("code is not exists")
	}
	return &giro, nil
}

func (e *repository) Update(id int, giro *model.Giro) (*model.Giro, error) {
	var upGiro = model.Giro{}
	err := e.DB.Table("giros").Where("id = ?", id).First(&upGiro).Update(&giro).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upGiro, nil
}

func (e *repository) Delete(id int) error {
	var giro = model.Giro{}
	err := e.DB.Table("giros").Where("id = ?", id).First(&giro).Delete(&giro).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}
