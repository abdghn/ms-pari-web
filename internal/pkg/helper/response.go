/*
 * Created on 01/04/22 17.20
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package helper

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponsePari struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    model.PariProduct `json:"data"`
}

type ResponsePariDetail struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    model.PariProductDetail `json:"data"`
}

type ResponsePaged struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	Total   int         `json:"total"`
}

type ProductResponse struct {
	*model.Product
	IsVerifiedByUser bool `json:"is_verified_by_user"`
}

type TransactionPreOrderResponse struct {
	*model.TransactionPreOrder
	IsVerifiedByUser bool `json:"is_verified_by_user"`
}

func HandleSuccess(c *gin.Context, data interface{}) {
	responseData := Response{
		Status:  "200",
		Message: "Success",
		Data:    data,
	}
	c.JSON(http.StatusOK, responseData)
}

func HandlePagedSuccess(c *gin.Context, data interface{}, page, size, total int) {
	responseData := ResponsePaged{
		Status:  "200",
		Message: "Success",
		Data:    data,
		Page:    page,
		Size:    size,
		Total:   total,
	}
	c.JSON(http.StatusOK, responseData)
}

func HandleError(c *gin.Context, status int, message string) {
	responseData := Response{
		Status:  strconv.Itoa(status),
		Message: message,
	}
	c.JSON(status, responseData)
}
