package product

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/product"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddProduct(c *gin.Context)
	ViewProductId(c *gin.Context)
	ViewProducts(c *gin.Context)
	ViewProductsBy(c *gin.Context)
	EditProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	SummaryProduct(c *gin.Context)
	VerificationProduct(c *gin.Context)
	PariProductTransaction(c *gin.Context)
}

type handler struct {
	usecase product.Usecase
}

func NewHandler(uc product.Usecase) Handler {
	return &handler{uc}
}

// AddProduct godoc
// @Summary Add new product
// @Schemes
// @Description add new product
// @Tags Product
// @Accept multipart/form-data
// @Param   file formData file false  "Upload Image"
// @Param        product  formData      request.Product  true  "Add product"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product [post]
func (e *handler) AddProduct(c *gin.Context) {
	var productModel request.Product

	err := c.ShouldBind(&productModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	path := "./internal/pkg/upload/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	file := productModel.File

	// generate new file name
	ext := filepath.Ext(file.Filename)
	currentTime := time.Now()
	filename := currentTime.Format("20060102150405") + ext

	tmpFile := path + filename
	if err = c.SaveUploadedFile(file, tmpFile); err != nil {
		helper.HandleError(c, http.StatusBadRequest, "failed to saving image")
		return
	}

	productModel.Image = "image/" + filename
	productModel.TmpImagePath = tmpFile

	if productModel.Name == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}

	newProduct, err := e.usecase.Create(&productModel)
	fmt.Println(err)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, newProduct)
}

// ViewProducts godoc
// @Summary Find All product
// @Schemes
// @Description find all product
// @Tags Product
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product [get]
func (e *handler) ViewProducts(c *gin.Context) {
	products, err := e.usecase.ReadAll()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	//if len(*products) == 0 {
	//	helper.HandleError(c, http.StatusNotFound, "list product is empty")
	//	return
	//}

	helper.HandleSuccess(c, products)
}

// ViewProductsBy godoc
// @Summary Find All product by Company ID
// @Schemes
// @Description find all product by company id
// @Param company_id path string true "Company ID"
// @Param   page     query    int     false        "Page"
// @Param   size      query    int     false        "Size"
// @Param   status      query    string     false        "Status"
// @Param   commodity      query    string     false        "Commodity"
// @Param   search      query    string     false        "Search"
// @Tags Product
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.ResponsePaged
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/company/{company_id} [get]
func (e *handler) ViewProductsBy(c *gin.Context) {
	var req request.ProductPaged
	var err error

	err = c.ShouldBindUri(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "company id has be number")
		return
	}

	req.CompanyID = companyID

	products, err := e.usecase.ReadAllBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	countProducts := e.usecase.Count(req)

	helper.HandlePagedSuccess(c, products, req.Page, req.Size, countProducts)
}

// ViewProductId FindProduct godoc
// @Summary Find product by id
// @Schemes
// @Description find product by id
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param   user_id      query    int     false        "User ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/{id} [get]
func (e *handler) ViewProductId(c *gin.Context) {
	var req request.ProductDetail
	var err error

	err = c.ShouldBindUri(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	err = c.ShouldBindQuery(&req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	productModel, err := e.usecase.ReadBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}

	helper.HandleSuccess(c, productModel)
}

// EditProduct UpdateProduct godoc
// @Summary update product by id
// @Schemes
// @Description update product by id
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param        product  body      request.Product  true  "Update product"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/{id} [put]
func (e *handler) EditProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	_, err = e.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	var tempProduct = model.Product{}
	err = c.Bind(&tempProduct)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempProduct.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}
	if tempProduct.Name == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	updatedProduct, err := e.usecase.Update(id, &tempProduct)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedProduct)
}

// DeleteProduct godoc
// @Summary Delete product by id
// @Schemes
// @Description delete product by id
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/{id} [delete]
func (e *handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	err = e.usecase.Delete(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, "success delete data")
}

// SummaryProduct godoc
// @Summary Find summary by Company ID
// @Schemes
// @Description find summary by company id
// @Tags Product
// @Accept  json
// @Produce  json
// @Param company_id path string true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/summary/{company_id} [get]
func (e *handler) SummaryProduct(c *gin.Context) {
	idStr := c.Param("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Print(err)
		helper.HandleError(c, http.StatusBadRequest, "company id has be number")
		return
	}
	productModel, err := e.usecase.Summary(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, productModel)
}

// VerificationProduct godoc
// @Summary Verification product
// @Schemes
// @Description verification product
// @Tags Product
// @Accept json
// @Produce json
// @Param        productUser  body      request.ProductUser  true  "Verification Product"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /product/verification [post]
func (e *handler) VerificationProduct(c *gin.Context) {
	var r = request.ProductUser{}
	err := c.Bind(&r)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	newProductUser, err := e.usecase.Verification(&r)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, newProductUser)
}

func (e *handler) PariProductTransaction(c *gin.Context) {
	var tempProduct = model.Product{}
	err := c.ShouldBind(&tempProduct)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	currentProduct, err := e.usecase.ReadByPariProductId(tempProduct.PariProductId)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}

	qty := currentProduct.Quantity - tempProduct.Quantity
	if qty < 0 {
		helper.HandleError(c, http.StatusInternalServerError, "stock")
		return
	}

	calculateProduct := model.Product{Quantity: qty}

	updatedProduct, err := e.usecase.Update(currentProduct.ID, &calculateProduct)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedProduct)
}
