package middleware

import (
	"fmt"
	"net/http"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"github.com/gin-gonic/gin"
)

//AuthorizeAPI -> to authorize Open API
func AuthorizeAPI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader("Authorization")
		if key == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization header found"})
		}

		if err := helper.ValidateApiKey(key); err != nil {

			fmt.Println("key", key, err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error()})

		}

	}

}
