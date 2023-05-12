/*
 * Created on 07/04/22 06.07
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package auth

import (
	"fmt"
	"net/http"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/request"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/auth"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(enforcer *casbin.Enforcer) gin.HandlerFunc
	BulkRegister(enforcer *casbin.Enforcer) gin.HandlerFunc
	Login(c *gin.Context)
	ValidateGiro(c *gin.Context)
	GetToken(c *gin.Context)
}

type handler struct {
	usecase auth.Usecase
}

func NewHandler(uc auth.Usecase) Handler {
	return &handler{uc}
}

// Register godoc
// @Summary Register
// @Schemes
// @Description register
// @Tags Auth
// @Accept json
// @Produce json
// @Param        user  body      request.User  true  "Register"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /register [post]
func (e *handler) Register(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user request.User
		err := c.Bind(&user)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
			return
		}

		if user.Email == "" {
			helper.HandleError(c, http.StatusBadRequest, "column cannot be empty")
			return
		}

		newUser, err := e.usecase.Register(user)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, err.Error())
			return
		}

		enforcer.AddGroupingPolicy(fmt.Sprint(newUser.ID), newUser.RoleName)
		newUser.Password = ""
		helper.HandleSuccess(c, newUser)
	}
}

// BulkRegister godoc
// @Summary Register
// @Schemes
// @Description register
// @Tags Auth
// @Accept json
// @Produce json
// @Param        users  body      request.Users  true  "Register"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /register/bulk [post]
func (e *handler) BulkRegister(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users request.Users
		err := c.Bind(&users)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
			return
		}

		newUsers, err := e.usecase.BulkRegister(users)
		if err != nil {
			helper.CommonLogger().Error(err)
			helper.HandleError(c, http.StatusInternalServerError, err.Error())
			return
		}

		for _, newUser := range newUsers {
			enforcer.AddGroupingPolicy(fmt.Sprint(newUser.ID), newUser.RoleName)
		}

		helper.HandleSuccess(c, newUsers)
	}
}

// Login godoc
// @Summary Login
// @Schemes
// @Description login
// @Tags Auth
// @Accept json
// @Produce json
// @Param        login  body      request.Login  true  "Login"
// @Success 201 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /login [post]
func (e *handler) Login(c *gin.Context) {
	var user = model.User{}
	err := c.Bind(&user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, "Oopss server someting wrong")
		return
	}

	dbUser, err := e.usecase.Login(&user)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	isTrue := helper.ComparePassword(dbUser.Password, user.Password)
	if !isTrue {
		helper.HandleError(c, http.StatusInternalServerError, "Password not matched")
		return
	}

	token := helper.GenerateToken(dbUser)
	result := map[string]interface{}{"token": token, "must_change_password": dbUser.MustChangePassword}
	helper.HandleSuccess(c, result)
}

// ValidateGiro godoc
// @Summary Find giro by code
// @Schemes
// @Description find giro by code
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param code path string true "Giro Code"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /validate_giro/{code} [get]
func (e *handler) ValidateGiro(c *gin.Context) {
	code := c.Param("code")

	r, err := e.usecase.ValidateGiro(code)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, r)
}

// GetToken godoc
// @Summary Get Token for Open API
// @Schemes
// @Description Get Token for Open API
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /token [get]
func (e *handler) GetToken(c *gin.Context) {
	clientKey := c.GetHeader("CLIENT_KEY")
	secretKey := c.GetHeader("SECRET_KEY")
	r, err := e.usecase.GetToken(clientKey, secretKey)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	helper.HandleSuccess(c, r)
}
