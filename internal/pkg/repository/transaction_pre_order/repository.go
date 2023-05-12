package transaction_pre_order

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(transactionPreOrder *model.TransactionPreOrder) (*model.TransactionPreOrder, error)
	ReadAll() (*[]model.TransactionPreOrder, error)
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.TransactionPreOrder, error)
	ReadById(id int) (*model.TransactionPreOrder, error)
	Update(id int, person *model.TransactionPreOrder) (*model.TransactionPreOrder, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(product *model.TransactionPreOrder) (*model.TransactionPreOrder, error) {
	err := e.DB.Save(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return product, nil
}

func (e *repository) ReadAll() (*[]model.TransactionPreOrder, error) {
	var transactionPreOrders []model.TransactionPreOrder
	err := e.DB.Select("*, p.name AS product_name, " +
		"p.image AS product_image, " +
		"p.commodity AS product_commodity, " +
		"p.min_price AS product_min_price, " +
		"p.max_price AS product_max_price, " +
		"p.product_created_at AS product_created_at," +
		"p.expired_at AS product_expired_at," +
		"p.is_pre_order AS product_is_pre_order," +
		"p.is_active AS product_is_active").
		Table("transaction_pre_orders").
		Joins("JOIN products p ON p.id = transaction_pre_orders.product_id").
		Find(&transactionPreOrders).Error

	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &transactionPreOrders, nil
}

func (e *repository) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.TransactionPreOrder, error) {
	var transactionPreOrders []model.TransactionPreOrder

	query := e.DB.Select("*, p.name AS product_name, " +
		"p.image AS product_image, " +
		"p.commodity AS product_commodity," +
		"p.min_price AS product_min_price," +
		"p.max_price AS product_max_price," +
		"p.product_created_at AS product_created_at," +
		"p.expired_at AS product_expired_at," +
		"p.is_pre_order AS product_is_pre_order," +
		"p.is_active AS product_is_active").
		Table("transaction_pre_orders").
		Joins("JOIN products p ON p.id = transaction_pre_orders.product_id").
		Where(criteria)

	if search != "" {
		query.Where("p.name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("transaction_pre_orders.created_at ASC").Limit(limit).Find(&transactionPreOrders).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &transactionPreOrders, nil
}

func (e *repository) ReadById(id int) (*model.TransactionPreOrder, error) {
	var transactionPreOrder = model.TransactionPreOrder{}
	err := e.DB.Select("*, p.name AS product_name, "+
		"p.image AS product_image, "+
		"p.commodity AS product_commodity, "+
		"p.min_price AS product_min_price, "+
		"p.max_price AS product_max_price, "+
		"p.product_created_at AS product_created_at,"+
		"p.expired_at AS product_expired_at,"+
		"p.is_pre_order AS product_is_pre_order,"+
		"p.is_active AS product_is_active").
		Table("transaction_pre_orders").
		Joins("JOIN products p ON p.id = transaction_pre_orders.product_id").
		Where("transaction_pre_orders.id = ?", id).First(&transactionPreOrder).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &transactionPreOrder, nil
}

func (e *repository) Update(id int, product *model.TransactionPreOrder) (*model.TransactionPreOrder, error) {
	var upTransactionPreOrder = model.TransactionPreOrder{}
	err := e.DB.Table("transaction_pre_orders").Where("id = ?", id).First(&upTransactionPreOrder).Update(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upTransactionPreOrder, nil
}

func (e *repository) Delete(id int) error {
	var product = model.TransactionPreOrder{}
	err := e.DB.Table("transaction_pre_orders").Where("id = ?", id).First(&product).Delete(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[transactionPreOrderRepository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *repository) Count(criteria map[string]interface{}) int {
	var result int
	err := e.DB.Table("transaction_pre_orders").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}
