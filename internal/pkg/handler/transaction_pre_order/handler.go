package transaction_pre_order

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	transactionPreOrder "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/transaction_pre_order"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddTransactionPreOrder(c *gin.Context)
	ViewTransactionPreOrderId(c *gin.Context)
	ViewTransactionPreOrders(c *gin.Context)
	ViewTransactionPreOrdersBy(c *gin.Context)
	EditTransactionPreOrder(c *gin.Context)
	DeleteTransactionPreOrder(c *gin.Context)
	SummaryTransactionPreOrder(c *gin.Context)
	VerificationTransactionPreOrder(c *gin.Context)
}

type handler struct {
	usecase transactionPreOrder.Usecase
}

func NewHandler(uc transactionPreOrder.Usecase) Handler {
	return &handler{uc}
}

// AddTransactionPreOrder godoc
// @Summary Add new transaction pre-order
// @Schemes
// @Description add new transaction pre-order
// @Tags Transaction PreOrder
// @Accept multipart/form-data
// @Param        transactionPreOrder  formData      request.TransactionPreOrder  true  "Add transaction pre-order"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder [post]
func (e *handler) AddTransactionPreOrder(c *gin.Context) {
	var transactionPreOrderModel request.TransactionPreOrder

	err := c.ShouldBind(&transactionPreOrderModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	newProduct, err := e.usecase.Create(&transactionPreOrderModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, newProduct)
}

// ViewTransactionPreOrders godoc
// @Summary Find All Transaction PreOrder
// @Schemes
// @Description find all transaction preorder
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder [get]
func (e *handler) ViewTransactionPreOrders(c *gin.Context) {
	transactionPreOrders, err := e.usecase.ReadAll()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, transactionPreOrders)
}

// ViewTransactionPreOrdersBy godoc
// @Summary Find All transaction preorder by Company ID
// @Schemes
// @Description find all transaction preorder by company id
// @Param company_id path string true "Company ID"
// @Param   page     query    int     false        "Page"
// @Param   size      query    int     false        "Size"
// @Param   status      query    string     false        "Status"
// @Param   commodity      query    string     false        "Commodity"
// @Param   search      query    string     false        "Search"
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.ResponsePaged
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/company/{company_id} [get]
func (e *handler) ViewTransactionPreOrdersBy(c *gin.Context) {
	var req request.TransactionPreOrderPaged
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

	transactionPreOrders, err := e.usecase.ReadAllBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	countTransactionPreOrders := e.usecase.Count(req)

	helper.HandlePagedSuccess(c, transactionPreOrders, req.Page, req.Size, countTransactionPreOrders)
}

// ViewTransactionPreOrderId FindTransactionPreOrder godoc
// @Summary Find transaction pre-order by id
// @Schemes
// @Description find transaction pre-order by id
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Param id path string true "Transaction PreOrder ID"
// @Param   user_id      query    int     false        "User ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/{id} [get]
func (e *handler) ViewTransactionPreOrderId(c *gin.Context) {
	var req request.TransactionPreOrderDetail
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

	transactionPreOrderModel, err := e.usecase.ReadBy(req)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}

	helper.HandleSuccess(c, transactionPreOrderModel)
}

// EditTransactionPreOrder UpdateTransactionPreOrder godoc
// @Summary update transaction pre-order by id
// @Schemes
// @Description update transaction pre-order by id
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Param id path string true "Transaction PreOrder ID"
// @Param        product  body      request.TransactionPreOrder  true  "Update Transaction PreOrder"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/{id} [put]
func (e *handler) EditTransactionPreOrder(c *gin.Context) {
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
	var tempTransactionPreOrder = model.TransactionPreOrder{}
	err = c.Bind(&tempTransactionPreOrder)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempTransactionPreOrder.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}
	updatedTransactionPreOrder, err := e.usecase.Update(id, &tempTransactionPreOrder)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedTransactionPreOrder)
}

// DeleteTransactionPreOrder godoc
// @Summary Delete transaction pre-order by id
// @Schemes
// @Description delete transaction pre-order by id
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Param id path string true "Transaction PreOrder ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/{id} [delete]
func (e *handler) DeleteTransactionPreOrder(c *gin.Context) {
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

// SummaryTransactionPreOrder godoc
// @Summary Find summary by Company ID
// @Schemes
// @Description find summary by company id
// @Tags Transaction PreOrder
// @Accept  json
// @Produce  json
// @Param company_id path string true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/summary/{company_id} [get]
func (e *handler) SummaryTransactionPreOrder(c *gin.Context) {
	idStr := c.Param("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Print(err)
		helper.HandleError(c, http.StatusBadRequest, "company id has be number")
		return
	}
	transactionPreOrderModel, err := e.usecase.Summary(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, transactionPreOrderModel)
}

// VerificationTransactionPreOrder godoc
// @Summary Verification transaction pre-order
// @Schemes
// @Description verification transaction pre-order
// @Tags Transaction PreOrder
// @Accept json
// @Produce json
// @Param        transactionPreOrderUser  body      request.TransactionPreOrderUser  true  "Verification Transaction PreOrder"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /transaction/preorder/verification [post]
func (e *handler) VerificationTransactionPreOrder(c *gin.Context) {
	var r = request.TransactionPreOrderUser{}
	err := c.Bind(&r)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	newTransactionPreOrderUser, err := e.usecase.Verification(&r)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.HandleSuccess(c, newTransactionPreOrderUser)
}
