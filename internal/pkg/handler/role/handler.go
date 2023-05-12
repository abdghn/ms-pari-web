package role

import (
	"net/http"
	"strconv"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/role"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddRole(enforcer *casbin.Enforcer) gin.HandlerFunc
	ViewRoleId(c *gin.Context)
	ViewRoles(c *gin.Context)
	EditRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

type handler struct {
	usecase role.Usecase
}

func NewHandler(uc role.Usecase) Handler {
	return &handler{uc}
}

// AddRole godoc
// @Summary Add new role
// @Schemes
// @Description add new role
// @Tags Role
// @Accept json
// @Produce json
// @Param        role  body      request.Role  true  "Add role"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /role [post]
func (e *handler) AddRole(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = model.Role{}
		err := c.Bind(&r)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
			return
		}
		if r.ID != 0 {
			helper.HandleError(c, http.StatusBadRequest, "input not permitted")
			return
		}

		if r.Name == "" {
			helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
			return
		}
		newRole, err := e.usecase.Create(&r)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, err.Error())
			return
		}

		//add policy
		if hasPolicy := enforcer.HasPolicy(newRole.Name, "report", "read"); !hasPolicy {
			enforcer.AddPolicy(newRole.Name, "report", "read")
		}
		if hasPolicy := enforcer.HasPolicy(newRole.Name, "report", "write"); !hasPolicy {
			enforcer.AddPolicy(newRole.Name, "report", "write")
		}

		helper.HandleSuccess(c, newRole)
	}
}

// ViewRoles godoc
// @Summary Find All role
// @Schemes
// @Description find all role
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /role [get]
func (e *handler) ViewRoles(c *gin.Context) {
	roles, err := e.usecase.ReadAll()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(*roles) == 0 {
		helper.HandleError(c, http.StatusNotFound, "list role is empty")
		return
	}
	helper.HandleSuccess(c, roles)
}

// ViewRoleId godoc
// @Summary Find role by id
// @Schemes
// @Description find role by id
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path string true "Role ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /role/{id} [get]
func (e *handler) ViewRoleId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	r, err := e.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, r)
}

// EditRole godoc
// @Summary update role by id
// @Schemes
// @Description update role by id
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path string true "Role ID"
// @Param        role  body      request.Role  true  "Update role"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /role/{id} [put]
func (e *handler) EditRole(c *gin.Context) {
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
	var tempRole = model.Role{}
	err = c.Bind(&tempRole)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempRole.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}
	if tempRole.Name == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	updatedRole, err := e.usecase.Update(id, &tempRole)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, updatedRole)
}

// DeleteRole godoc
// @Summary Delete role by id
// @Schemes
// @Description delete role by id
// @Tags Role
// @Accept  json
// @Produce  json
// @Param id path string true "Role ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /role/{id} [delete]
func (e *handler) DeleteRole(c *gin.Context) {
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
