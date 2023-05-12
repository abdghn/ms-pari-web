/*
 * Created on 01/04/22 15.31
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package user

import (
	"net/http"
	"strconv"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/user"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddUser(c *gin.Context)
	ViewUserId(c *gin.Context)
	ViewUsers(c *gin.Context)
	EditUser(c *gin.Context)
	ChangePassword(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type handler struct {
	usecase user.Usecase
}

func NewHandler(uc user.Usecase) Handler {
	return &handler{uc}
}

// AddUser godoc
// @Summary Add new user
// @Schemes
// @Description add new user
// @Tags User
// @Accept json
// @Produce json
// @Param        user  body      request.User  true  "Add user"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user [post]
func (e *handler) AddUser(c *gin.Context) {
	var userModel = model.User{}
	err := c.Bind(&userModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if userModel.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}

	if userModel.Name == "" || userModel.Email == "" || userModel.Password == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	newUser, err := e.usecase.Create(&userModel)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, newUser)
}

// ViewUsers godoc
// @Summary Find All user
// @Schemes
// @Description find all user
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user [get]
func (e *handler) ViewUsers(c *gin.Context) {
	users, err := e.usecase.ReadAll()
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(*users) == 0 {
		helper.HandleError(c, http.StatusNotFound, "list user is empty")
		return
	}
	helper.HandleSuccess(c, users)
}

// ViewUserId godoc
// @Summary Find user by id
// @Schemes
// @Description find user by id
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user/{id} [get]
func (e *handler) ViewUserId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}
	u, err := e.usecase.ReadById(id)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, u)
}

// EditUser UpdateUser godoc
// @Summary update user by id
// @Schemes
// @Description update user by id
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param        user  body      request.User  true  "Update user"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user/{id} [put]
func (e *handler) EditUser(c *gin.Context) {
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
	var tempUser = model.User{}
	err = c.Bind(&tempUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}
	if tempUser.ID != 0 {
		helper.HandleError(c, http.StatusBadRequest, "input not permitted")
		return
	}
	if tempUser.Name == "" || tempUser.Email == "" || tempUser.Password == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}
	u, err := e.usecase.Update(id, &tempUser)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.HandleSuccess(c, u)
}

// DeleteUser godoc
// @Summary Delete user by id
// @Schemes
// @Description delete user by id
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user/{id} [delete]
func (e *handler) DeleteUser(c *gin.Context) {
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

// ChangePassword godoc
// @Summary Add new user
// @Schemes
// @Description add new user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param        password  body      request.ChangePassword  true  "Change Password"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Security BearerAuth
// @Router /user/change_password/{id} [put]
func (e *handler) ChangePassword(c *gin.Context) {
	var changePassword = request.ChangePassword{}
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusBadRequest, "id has be number")
		return
	}

	err = c.Bind(&changePassword)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	if changePassword.Password == "" {
		helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
		return
	}

	changePassword.UserID = userID

	helper.HashPassword(&changePassword.Password)
	m, err := e.usecase.ChangePassword(changePassword)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	token := helper.GenerateToken(m)
	result := map[string]interface{}{"token": token, "must_change_password": m.MustChangePassword}
	helper.HandleSuccess(c, result)
}
