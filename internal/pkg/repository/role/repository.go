package role

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(person *model.Role) (*model.Role, error)
	ReadAll() (*[]model.Role, error)
	ReadById(id int) (*model.Role, error)
	ReadByName(name string) (*model.Role, error)
	Update(id int, person *model.Role) (*model.Role, error)
	Delete(id int) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(role *model.Role) (*model.Role, error) {
	err := e.DB.Save(&role).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return role, nil
}

func (e *repository) ReadAll() (*[]model.Role, error) {
	var roles []model.Role
	err := e.DB.Find(&roles).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &roles, nil
}

func (e *repository) ReadById(id int) (*model.Role, error) {
	var role = model.Role{}
	err := e.DB.Table("roles").Where("id = ?", id).First(&role).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[roleRepository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &role, nil
}

func (e *repository) ReadByName(name string) (*model.Role, error) {
	var role = model.Role{}
	err := e.DB.Table("roles").Where("name = ?", name).First(&role).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("name is not exists")
	}
	return &role, nil
}

func (e *repository) Update(id int, role *model.Role) (*model.Role, error) {
	var upRole = model.Role{}
	err := e.DB.Table("roles").Where("id = ?", id).First(&upRole).Update(&role).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upRole, nil
}

func (e *repository) Delete(id int) error {
	var role = model.Role{}
	err := e.DB.Table("roles").Where("id = ?", id).First(&role).Delete(&role).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[repository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}
