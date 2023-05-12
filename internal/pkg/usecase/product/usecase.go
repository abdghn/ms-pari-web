package product

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/product"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/product_user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/role"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Usecase interface {
	Create(product *request.Product) (*model.Product, error)
	ReadAll() (*[]model.Product, error)
	ReadAllBy(req request.ProductPaged) (*[]model.Product, error)
	ReadById(id int) (*model.Product, error)
	ReadByPariProductId(pariProductId string) (*model.Product, error)
	ReadBy(req request.ProductDetail) (*helper.ProductResponse, error)
	Update(id int, product *model.Product) (*model.Product, error)
	Delete(id int) error
	Count(req request.ProductPaged) int
	Summary(companyId int) (interface{}, error)
	Verification(productUser *request.ProductUser) (*helper.ProductResponse, error)
}

type usecase struct {
	productRepository     product.Repository
	productUserRepository product_user.Repository
	userRepository        user.Repository
	roleRepository        role.Repository
}

func NewUsecase(productRepository product.Repository, productUserRepository product_user.Repository, userRepository user.Repository, roleRepository role.Repository) Usecase {
	return &usecase{productRepository, productUserRepository, userRepository, roleRepository}
}

func (e *usecase) Create(product *request.Product) (*model.Product, error) {
	//layout := "2006-01-02"
	//startDateString := product.ProductCreatedAt.Format(layout)
	//endDateString := product.ExpiredAt.Format(layout)
	//t, _ := time.Parse(layout, product.ProductCreatedAt)
	//t2, _ := time.Parse(layout, product.ExpiredAt)

	p := &model.Product{
		Image:            product.Image,
		Name:             product.Name,
		Description:      product.Description,
		Quantity:         product.Quantity,
		UnitQuantity:     product.UnitQuantity,
		Price:            product.Price,
		UnitPrice:        product.UnitPrice,
		Status:           product.Status,
		ProductCreatedAt: product.ProductCreatedAt,
		ExpiredAt:        product.ExpiredAt,
		Commodity:        product.Commodity,
		CompanyID:        product.CompanyID,
		IsPreOrder:       product.IsPreOrder,
		MinPrice:         product.MinPrice,
		MaxPrice:         product.MaxPrice,
		IsActive:         product.IsActive,
		TmpImagePath:     product.TmpImagePath,
	}

	return e.productRepository.Create(p)
}

func (e *usecase) ReadAll() (*[]model.Product, error) {
	return e.productRepository.ReadAll()
}

func (e *usecase) ReadAllBy(req request.ProductPaged) (*[]model.Product, error) {
	criteria := make(map[string]interface{})
	criteria["company_id"] = req.CompanyID

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	if req.Commodity != "" {
		criteria["commodity"] = req.Commodity
	}

	fmt.Println(req)

	return e.productRepository.ReadAllBy(criteria, req.Search, req.Page, req.Size)
}

func (e *usecase) ReadById(id int) (*model.Product, error) {
	return e.productRepository.ReadById(id)
}

func (e *usecase) ReadByPariProductId(pariProductId string) (*model.Product, error) {
	return e.productRepository.ReadByPariProductId(pariProductId)
}

func (e *usecase) ReadBy(request request.ProductDetail) (*helper.ProductResponse, error) {
	var isVerifiedByUser bool
	productModel, err := e.productRepository.ReadById(request.ID)

	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	if request.UserID != 0 {
		countProductUser := e.productUserRepository.Count(map[string]interface{}{"product_id": request.ID, "user_id": request.UserID})
		if countProductUser > 0 {
			isVerifiedByUser = true
		}
	}

	extraFields := map[string]string{
		"corporate_id": strconv.Itoa(productModel.CompanyID),
		"product_id":   productModel.PariProductId,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range extraFields {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", viper.Get("API_PARI_CORPORATE").(string)+enum.DetailProduct.String(), body)
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", viper.Get("API_KEY_PARI_CORPORATE").(string))

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("http.Do() error: %v\n", string(b))
		return nil, fmt.Errorf(string(b))
	}

	fmt.Printf("http.Do() success: %v\n", resp.Body)

	defer resp.Body.Close()

	var response helper.ResponsePariDetail

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	fmt.Println(response.Data)

	productModel.Name = response.Data.ProductName
	productModel.Image = response.Data.Images
	productModel.Transaction = response.Data.Transaction

	result := &helper.ProductResponse{Product: productModel, IsVerifiedByUser: isVerifiedByUser}

	return result, nil
}

func (e *usecase) Update(id int, product *model.Product) (*model.Product, error) {
	return e.productRepository.Update(id, product)
}

func (e *usecase) Delete(id int) error {
	return e.productRepository.Delete(id)
}

func (e *usecase) Verification(request *request.ProductUser) (*helper.ProductResponse, error) {

	productModel, err := e.productRepository.ReadById(request.ProductID)
	if err != nil {
		helper.CommonLogger().Error(err)
	}

	productUser, err := e.productUserRepository.ReadBy(map[string]interface{}{"product_id": request.ProductID, "user_id": request.UserID, "company_id": request.CompanyID})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("failed finding product user: %v", err)
	}

	r, err := e.roleRepository.ReadById(request.RoleID)
	if err != nil {
		helper.CommonLogger().Error(err)
		return nil, err
	}

	// checking whether productUser exists or not
	if productUser == nil {
		pu := &model.ProductUser{ProductID: request.ProductID, UserID: request.UserID, CompanyID: request.CompanyID}
		_, err := e.productUserRepository.Create(pu)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}
	}

	// checking product has been approved by all user in company
	countProductUser := e.productUserRepository.Count(map[string]interface{}{"company_id": request.CompanyID, "product_id": request.ProductID})
	countUser := e.userRepository.Count(map[string]interface{}{"company_id": request.CompanyID, "role_id": r.ID})
	if countUser == countProductUser {
		var isPreOrder int

		if productModel.IsPreOrder {
			isPreOrder = 1
		} else {
			isPreOrder = 0
		}

		file, err := os.Open(productModel.Image)
		if err != nil {
			return nil, err
		}

		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		file.Close()

		extraFields := map[string]string{
			"corporate_id":      strconv.Itoa(productModel.CompanyID),
			"product_name":      productModel.Name,
			"product_commodity": productModel.ProductCreatedAt,
			"date_production":   productModel.ProductCreatedAt,
			"expires_date":      productModel.ExpiredAt,
			"price":             strconv.Itoa(int(productModel.Price)),
			"minPrice":          strconv.Itoa(int(productModel.MinPrice)),
			"maxPrice":          strconv.Itoa(int(productModel.MaxPrice)),
			"isPreOrder":        strconv.Itoa(isPreOrder),
			"status":            strconv.Itoa(1),
			"description":       productModel.Description,
			"quantity":          strconv.Itoa(productModel.Quantity),
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("images", fi.Name())
		if err != nil {
			return nil, err
		}

		_, err = part.Write(fileContents)
		if err != nil {
			return nil, err
		}

		for key, val := range extraFields {
			_ = writer.WriteField(key, val)
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", viper.Get("API_PARI_CORPORATE").(string)+enum.CreateProduct.String(), body)
		if err != nil {
			fmt.Printf("http.NewRequest() error: %v\n", err)
			return nil, err
		}

		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Add("Authorization", viper.Get("API_KEY_PARI_CORPORATE").(string))

		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			fmt.Printf("http.Do() error: %v\n", err)
			return nil, err
		}

		if resp.StatusCode != 200 {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("http.Do() error: %v\n", string(b))
			return nil, fmt.Errorf(string(b))
		}

		defer resp.Body.Close()

		var response helper.ResponsePari

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err
		}

		fmt.Println(response.Data)

		productModel.Status = enum.Approved
		productModel.PariProductId = response.Data.ID
		productModel.Image = ""
		productModel.TmpImagePath = ""

		_, err = e.productRepository.Update(productModel.ID, productModel)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}

		err = os.Remove(productModel.TmpImagePath)
		if err != nil {
			helper.CommonLogger().Error(err)
			return nil, err
		}
	}

	result := &helper.ProductResponse{Product: productModel, IsVerifiedByUser: true}
	return result, nil
}

func (e *usecase) Count(req request.ProductPaged) int {
	criteria := make(map[string]interface{})
	criteria["company_id"] = req.CompanyID

	if req.Status != "" {
		criteria["status"] = req.Status
	}

	if req.Commodity != "" {
		criteria["commodity"] = req.Commodity
	}

	return e.productRepository.Count(criteria)
}

func (e *usecase) Summary(companyId int) (interface{}, error) {

	allProduct := e.productRepository.Count(map[string]interface{}{"company_id": companyId})
	processingProduct := e.productRepository.Count(map[string]interface{}{"company_id": companyId, "status": "processing"})
	approvedProduct := e.productRepository.Count(map[string]interface{}{"company_id": companyId, "status": "approved"})
	rejectedProduct := e.productRepository.Count(map[string]interface{}{"company_id": companyId, "status": "rejected"})

	return map[string]interface{}{
		"all_product":        allProduct,
		"processing_product": processingProduct,
		"approved_product":   approvedProduct,
		"rejected_product":   rejectedProduct,
	}, nil
}
