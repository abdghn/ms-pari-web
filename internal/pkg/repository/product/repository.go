package product

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	Create(person *model.Product) (*model.Product, error)
	ReadAll() (*[]model.Product, error)
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Product, error)
	ReadById(id int) (*model.Product, error)
	ReadByPariProductId(pariProductId string) (*model.Product, error)
	Update(id int, person *model.Product) (*model.Product, error)
	Delete(id int) error
	Count(criteria map[string]interface{}) int
	CreatePariProduct(product *model.Product) (*model.Product, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) Repository {
	return &repository{DB}
}

func (e *repository) Create(product *model.Product) (*model.Product, error) {
	err := e.DB.Save(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return product, nil
}

func (e *repository) ReadAll() (*[]model.Product, error) {
	var products []model.Product
	err := e.DB.Find(&products).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.ReadAll] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &products, nil
}

func (e *repository) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Product, error) {
	var products []model.Product

	query := e.DB.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at ASC").Limit(limit).Find(&products).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &products, nil
}

func (e *repository) ReadById(id int) (*model.Product, error) {
	var product = model.Product{}
	err := e.DB.Table("products").Where("id = ?", id).First(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &product, nil
}

func (e *repository) ReadByPariProductId(pariProductId string) (*model.Product, error) {
	var product = model.Product{}
	if err := e.DB.Table("products").Where("pari_product_id = ?", pariProductId).First(&product).Error; err != nil {
		fmt.Printf("[productRepository.ReadByPariProductId] error execute query %v \n", err)
		return nil, fmt.Errorf("pari product id is not exists")
	}
	return &product, nil
}

func (e *repository) Update(id int, product *model.Product) (*model.Product, error) {
	var upProduct = model.Product{}
	err := e.DB.Table("products").Where("id = ?", id).First(&upProduct).Update(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upProduct, nil
}

func (e *repository) Delete(id int) error {
	var product = model.Product{}
	err := e.DB.Table("products").Where("id = ?", id).First(&product).Delete(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (e *repository) Count(criteria map[string]interface{}) int {
	var result int
	err := e.DB.Table("products").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (e *repository) CreatePariProduct(product *model.Product) (*model.Product, error) {
	err := e.DB.Save(&product).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[productRepository.CreatePariProduct] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}
	return product, nil
}
