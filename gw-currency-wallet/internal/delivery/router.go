package delivery

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetBalance(c *gin.Context)
	Deposit(c *gin.Context)
	Withdraw(c *gin.Context)
	Refresh(c *gin.Context)
	ExchangeRatesHandler(c *gin.Context)
	ExchangeHandler(c *gin.Context)
}

func NewRouter(router *gin.Engine, authMiddleware gin.HandlerFunc, c Controller) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRoutes := router.Group("/api/v1")
	{
		publicRoutes.POST("/register", c.Register)
		publicRoutes.POST("/login", c.Login)
		publicRoutes.POST("/refresh", c.Refresh)
	}

	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(authMiddleware)
	walletRoutes := protectedRoutes.Group("/wallet")
	{

		walletRoutes.GET("/balance", c.GetBalance)
		walletRoutes.POST("/deposit", c.Deposit)
		walletRoutes.POST("/withdraw", c.Withdraw)
	}
	protectedRoutes.GET("/exchange/rates", c.ExchangeRatesHandler)
	protectedRoutes.POST("/exchange", c.ExchangeHandler)
}
