/*
 * Created on 01/04/22 14.02
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package main

import (
	"net/http"
	"time"

	transactionPreOrderHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/transaction_pre_order"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/helper"
	transactionPreOrderRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/transaction_pre_order"
	transactionPreOrderUserRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/transaction_pre_order_user"
	transactionPreOrderUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/transaction_pre_order"

	productUserRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/product_user"

	"bitbucket.org/bridce/ms-pari-web/docs"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/config"
	authHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/auth"
	companyHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/company"
	productHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/product"
	roleHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/role"
	userHandler "bitbucket.org/bridce/ms-pari-web/internal/pkg/handler/user"
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/middleware"
	companyRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/company"
	giroRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/giro"
	productRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/product"
	roleRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/role"
	userRepository "bitbucket.org/bridce/ms-pari-web/internal/pkg/repository/user"
	authUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/auth"
	companyUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/company"
	productUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/product"
	roleUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/role"
	userUsecase "bitbucket.org/bridce/ms-pari-web/internal/pkg/usecase/user"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PARI Korporat
// @version 1.0
// @description PARI Korporat REST API
// @BasePath  /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
func main() {

	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		helper.CommonLogger().Error(err)
		return
	}

	port := viper.Get("PORT").(string)
	dbUser := viper.Get("DB_USER").(string)
	dbPass := viper.Get("DB_PASSWORD").(string)
	dbHost := viper.Get("DB_HOST").(string)
	dbPort := viper.Get("DB_PORT").(string)
	dbName := viper.Get("DB_NAME").(string)

	db := config.DbConnect(dbUser, dbPass, dbHost, dbPort, dbName)
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			helper.CommonLogger().Error(err)
			return
		}
	}(db)

	// Initialize  casbin adapter
	adapter := gormadapter.NewAdapterByDB(db)

	// Load model configuration file and policy store adapter
	enforcer := casbin.NewEnforcer("./internal/pkg/config/rbac_model.conf", adapter)

	//add policy
	if hasPolicy := enforcer.HasPolicy("superadmin", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("supeardmin", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("verificator", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("verificator", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("user", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("user", "report", "read")
	}

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	err = router.SetTrustedProxies(nil)
	if err != nil {
		helper.CommonLogger().Error(err)
		return
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3003", viper.Get("ALLOW_ORIGIN").(string)},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// init repositories
	userRepo := userRepository.NewRepository(db)
	roleRepo := roleRepository.NewRepository(db)
	companyRepo := companyRepository.NewRepository(db)
	giroRepo := giroRepository.NewRepository(db)
	productRepo := productRepository.NewRepository(db)
	productUserRepo := productUserRepository.NewRepository(db)
	transactionPreOrderRepo := transactionPreOrderRepository.NewRepository(db)
	transactionPreOrderUserRepo := transactionPreOrderUserRepository.NewRepository(db)

	// init usecases
	userUC := userUsecase.NewUsecase(userRepo)
	authUC := authUsecase.NewUsecase(userRepo, giroRepo, roleRepo, companyRepo)
	roleUC := roleUsecase.NewUsecase(roleRepo)
	companyUC := companyUsecase.NewUsecase(companyRepo)
	productUC := productUsecase.NewUsecase(productRepo, productUserRepo, userRepo, roleRepo)
	transactionPreOrderUC := transactionPreOrderUsecase.NewUsecase(transactionPreOrderRepo, transactionPreOrderUserRepo, userRepo, roleRepo)

	// init handlers
	userH := userHandler.NewHandler(userUC)
	authH := authHandler.NewHandler(authUC)
	roleH := roleHandler.NewHandler(roleUC)
	companyH := companyHandler.NewHandler(companyUC)
	productH := productHandler.NewHandler(productUC)
	transactionPreOrderH := transactionPreOrderHandler.NewHandler(transactionPreOrderUC)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", authH.Register(enforcer))
		v1.POST("/register/bulk", authH.BulkRegister(enforcer))
		v1.POST("/login", authH.Login)
		v1.GET("/validate_giro/:code", authH.ValidateGiro)

		// init open api
		v1.GET("/token", authH.GetToken)
		v1.GET("/company/:id", middleware.AuthorizeAPI(), companyH.ViewCompanyId)
		v1.POST("/product/preorder", middleware.AuthorizeAPI(), transactionPreOrderH.AddTransactionPreOrder)
		v1.POST("/product/transaction", middleware.AuthorizeAPI(), productH.PariProductTransaction)

		// init user routes
		user := v1.Group("/user", middleware.AuthorizeJWT())
		{
			user.GET("", middleware.Authorize("report", "read", enforcer), userH.ViewUsers)
			user.POST("", middleware.Authorize("report", "write", enforcer), userH.AddUser)
			user.PUT("/change_password/:id", userH.ChangePassword)
			user.GET("/:id", middleware.Authorize("report", "read", enforcer), userH.ViewUserId)
			user.PUT("/:id", middleware.Authorize("report", "write", enforcer), userH.EditUser)
			user.DELETE("/:id", middleware.Authorize("report", "write", enforcer), userH.DeleteUser)
		}

		// init role routes
		role := v1.Group("/role")
		{
			role.GET("", roleH.ViewRoles)
			role.POST("", roleH.AddRole(enforcer))
			role.GET("/:id", roleH.ViewRoleId)
			role.PUT("/:id", roleH.EditRole)
			role.DELETE("/:id", roleH.DeleteRole)
		}

		// init company routes
		company := v1.Group("/company")
		{
			company.GET("", companyH.ViewCompanies)
			company.POST("", companyH.AddCompany)
			company.PUT("/:id", companyH.EditCompany)
			company.DELETE("/:id", companyH.DeleteCompany)
		}

		// init product routes
		product := v1.Group("/product", middleware.AuthorizeJWT())
		{
			product.GET("", productH.ViewProducts)
			product.GET("/company/:company_id", productH.ViewProductsBy)
			product.POST("", productH.AddProduct)
			product.GET("/:id", productH.ViewProductId)
			product.GET("/summary/:company_id", productH.SummaryProduct)
			product.PUT("/:id", productH.EditProduct)
			product.DELETE("/:id", productH.DeleteProduct)
			product.POST("/verification", productH.VerificationProduct)
		}

		// init transaction pre order routes
		tpo := v1.Group("/transaction/preorder", middleware.AuthorizeJWT())
		{
			tpo.GET("", transactionPreOrderH.ViewTransactionPreOrders)
			tpo.GET("/company/:company_id", transactionPreOrderH.ViewTransactionPreOrdersBy)
			tpo.POST("", transactionPreOrderH.AddTransactionPreOrder)
			tpo.GET("/:id", transactionPreOrderH.ViewTransactionPreOrderId)
			tpo.GET("/summary/:company_id", transactionPreOrderH.SummaryTransactionPreOrder)
			tpo.PUT("/:id", transactionPreOrderH.EditTransactionPreOrder)
			tpo.DELETE("/:id", transactionPreOrderH.DeleteTransactionPreOrder)
			tpo.POST("/verification", transactionPreOrderH.VerificationTransactionPreOrder)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.StaticFS("/image", http.Dir("internal/pkg/upload"))

	err = router.Run(":" + port)
	if err != nil {
		helper.CommonLogger().Error(err)
		helper.CommonLogger().Error(err)
	}
}
