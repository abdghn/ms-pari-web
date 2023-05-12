package middleware

import (
	"fmt"

	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	"github.com/casbin/casbin"

	"github.com/gin-gonic/gin"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		sub, existed := c.Get("userID")
		if !existed {
			c.AbortWithStatusJSON(401, gin.H{"msg": "User hasn't logged in yet"})
			return
		}

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			helper.CommonLogger().Error(err)
			c.AbortWithStatusJSON(500, gin.H{"msg": "Failed to load policy from DB"})
			return
		}

		// Casbin enforces policy
		ok := enforcer.Enforce(fmt.Sprint(sub), obj, act)

		//if err != nil {
		helper.CommonLogger().Error(err)
		//	c.AbortWithStatusJSON(500, gin.H{"msg": "Error occurred when authorizing user"})
		//	return
		//}

		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"msg": "You are not authorized"})
			return
		}
		c.Next()
	}
}
