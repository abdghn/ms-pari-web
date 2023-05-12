package company

import (
	"net/http"
	"strconv"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/company"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddCompany(c *gin.Context)
	ViewCompanyId(c *gin.Context)
	ViewCompanies(c *gin.Context)
	EditCompany(c *gin.Context)
	DeleteCompany(c *gin.Context)
}

type handler struct {
	usecase company.Usecase
}

func NewHandler(uc company.Usecase) Handler {
	return &handler{uc}
}

// AddCompany godoc
// @Summary Add new company
// @Schemes
// @Description add new company
// @Tags Company
// @Accept json
// @Produce json
// @Param        company  body      request.Company  true  "Add company"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /company [post]
func (e *handler) AddCompany(c *gin.Context) {
	var companyModel = model.Company{}
	err := c.Bind(&companyModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if companyModel.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}

	if companyModel.Name == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	newCompany, err := e.usecase.Create(&companyModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, newCompany)
}

// ViewCompanies godoc
// @Summary Find All company
// @Schemes
// @Description find all company
// @Tags Company
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /company [get]
func (e *handler) ViewCompanies(c *gin.Context) {
	companys, err := e.usecase.ReadAll()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(*companys) == 0 {
		helper.HandleError(c, http.StatusNotFound, "list company is empty")
		return
	}
	helper.HandleSuccess(c, companys)
}

// ViewCompanyId FindCompany godoc
// @Summary Find company by id
// @Schemes
// @Description find company by id
// @Tags Company
// @Accept  json
// @Produce  json
// @Param id path string true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /company/{id} [get]
func (e *handler) ViewCompanyId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	companyModel, err := e.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, companyModel)
}

// EditCompany UpdateCompany godoc
// @Summary update company by id
// @Schemes
// @Description update company by id
// @Tags Company
// @Accept  json
// @Produce  json
// @Param id path string true "Company ID"
// @Param        company  body      request.Company  true  "Update company"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /company/{id} [put]
func (e *handler) EditCompany(c *gin.Context) {
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
	var tempCompany = model.Company{}
	err = c.Bind(&tempCompany)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempCompany.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}
	if tempCompany.Name == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	updatedCompany, err := e.usecase.Update(id, &tempCompany)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedCompany)
}

// DeleteCompany godoc
// @Summary Delete company by id
// @Schemes
// @Description delete company by id
// @Tags Company
// @Accept  json
// @Produce  json
// @Param id path string true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /company/{id} [delete]
func (e *handler) DeleteCompany(c *gin.Context) {
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
