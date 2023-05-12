package company

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(person *model.Company) (*model.Company, error)
	ReadAll() (*[]model.Company, error)
	ReadById(id int) (*model.Company, error)
	Update(id int, person *model.Company) (*model.Company, error)
	Delete(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(company *model.Company) (*model.Company, error) {
	err := e.DB.Save(&company).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return company, nil
}

func (e *repository) ReadAll() (*[]model.Company, error) {
	var companies []model.Company
	err := e.DB.Find(&companies).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &companies, nil
}

func (e *repository) ReadById(id int) (*model.Company, error) {
	var company = model.Company{}
	err := e.DB.Table("companies").Where("id = ?", id).First(&company).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &company, nil
}

func (e *repository) Update(id int, company *model.Company) (*model.Company, error) {
	var upCompany = model.Company{}
	err := e.DB.Table("companies").Where("id = ?", id).First(&upCompany).Update(&company).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upCompany, nil
}

func (e *repository) Delete(id int) error {
	var company = model.Company{}
	err := e.DB.Table("companies").Where("id = ?", id).First(&company).Delete(&company).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}
