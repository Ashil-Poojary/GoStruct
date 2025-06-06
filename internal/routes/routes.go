package routes

import (
	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/db"
	"github.com/ashil-poojary/gostruct/internal/handlers"
	"github.com/ashil-poojary/gostruct/internal/middleware"
	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	if db.DB1 == nil || db.DB2 == nil {
		panic("Databases are not initialized")
	}

	userRepo := repository.NewUserRepository(db.DB1, db.DB2)
	refreshTokenRepo := repository.NewAuthRepo(db.DB1)
	orderRepo := repository.NewOrderRepository(db.DB1, db.DB2)

	authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	// Get Microsoft OAuth handlers
	msLogin, msCallback := handlers.NewMicrosoftOAuthHandler(cfg.MicrosoftOAuth)

	v1 := r.Group("/v1")
	{
		public := v1.Group("/auth")
		{
			public.POST("/login", authHandler.Login)
			public.POST("/refresh", authHandler.RefreshToken)
			public.POST("/logout", authHandler.Logout)
			public.GET("/microsoft/login", msLogin)
			public.GET("/microsoft/callback", msCallback)
		}

		protected := v1.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		orders := protected.Group("/orders")
		{
			orders.GET("/", orderHandler.GetAllOrders)
		}

		protectedAdmin := protected.Group("/admin")
		protectedAdmin.Use(middleware.RoleAuthMiddleware(utils.Constants.RoleAdmin))
		{
			protectedAdmin.GET("/users", userHandler.GetAllUsers)
			protectedAdmin.GET("/users/:id", userHandler.GetUserByID)
		}
	}

	return r
}
