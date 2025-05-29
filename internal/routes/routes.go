package routes

import (
	"github.com/ashil-poojary/gostruct/internal/db"
	"github.com/ashil-poojary/gostruct/internal/handlers"
	"github.com/ashil-poojary/gostruct/internal/middleware"
	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Assuming db.DB1 and db.DB2 are initialized in the db package
	// and are accessible here.

	//check if db.DB1 and db.DB2 are initialized
	if db.DB1 == nil || db.DB2 == nil {
		panic("Databases are not initialized. Please check your database connection settings.")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB1, db.DB2)
	refreshTokenRepo := repository.NewAuthRepo(db.DB1)
	orderRepo := repository.NewOrderRepository(db.DB1, db.DB2)
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	v1 := r.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)

		}
		// protected routes
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		orders := protected.Group("/orders")
		{
			orders.GET("/", orderHandler.GetAllOrders)
		}
		// Only allow admins
		protectedAdmin := protected.Group("/admin")
		protectedAdmin.Use(middleware.RoleAuthMiddleware(utils.Constants.RoleAdmin))
		{
			protectedAdmin.GET("/users", userHandler.GetAllUsers)
			protectedAdmin.GET("/users/:id", userHandler.GetUserByID)
		}

	}

	return r
}
