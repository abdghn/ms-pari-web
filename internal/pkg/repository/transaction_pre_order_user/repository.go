package transaction_pre_order_user

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(transactionPreOrderUser *model.TransactionPreOrderUser) (*model.TransactionPreOrderUser, error)
	ReadAll() (*[]model.TransactionPreOrderUser, error)
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.TransactionPreOrderUser, error)
	ReadById(id int) (*model.TransactionPreOrderUser, error)
	ReadBy(criteria map[string]interface{}) (*model.TransactionPreOrderUser, error)
	Update(id int, person *model.TransactionPreOrderUser) (*model.TransactionPreOrderUser, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(transactionPreOrderUser *model.TransactionPreOrderUser) (*model.TransactionPreOrderUser, error) {
	err := e.DB.Save(&transactionPreOrderUser).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return transactionPreOrderUser, nil
}

func (e *repository) ReadAll() (*[]model.TransactionPreOrderUser, error) {
	var transactionPreOrderUsers []model.TransactionPreOrderUser
	err := e.DB.Find(&transactionPreOrderUsers).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &transactionPreOrderUsers, nil
}

func (e *repository) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.TransactionPreOrderUser, error) {
	var transactionPreOrderUsers []model.TransactionPreOrderUser

	query := e.DB.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at DESC").Limit(limit).Find(&transactionPreOrderUsers).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &transactionPreOrderUsers, nil
}

func (e *repository) ReadById(id int) (*model.TransactionPreOrderUser, error) {
	var transactionPreOrderUser = model.TransactionPreOrderUser{}
	err := e.DB.Table("transaction_pre_order_users").Where("id = ?", id).First(&transactionPreOrderUser).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &transactionPreOrderUser, nil
}

func (e *repository) Update(id int, transactionPreOrderUser *model.TransactionPreOrderUser) (*model.TransactionPreOrderUser, error) {
	var upTransactionPreOrderUser = model.TransactionPreOrderUser{}
	err := e.DB.Table("transaction_pre_order_users").Where("id = ?", id).First(&upTransactionPreOrderUser).Update(&transactionPreOrderUser).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upTransactionPreOrderUser, nil
}

func (e *repository) Delete(id int) error {
	var transactionPreOrderUser = model.TransactionPreOrderUser{}
	err := e.DB.Table("transaction_pre_order_users").Where("id = ?", id).First(&transactionPreOrderUser).Delete(&transactionPreOrderUser).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *repository) Count(criteria map[string]interface{}) int {
	var result int
	err := e.DB.Table("transaction_pre_order_users").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (e *repository) ReadBy(criteria map[string]interface{}) (*model.TransactionPreOrderUser, error) {
	var transactionPreOrderUser = model.TransactionPreOrderUser{}
	err := e.DB.Table("transaction_pre_order_users").Where(criteria).First(&transactionPreOrderUser).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderUserRepository.ReadBy] error execute query %v \n", err)
		return nil, err
	}
	return &transactionPreOrderUser, nil
}
