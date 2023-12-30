package router

import (
	"vix-btpns/controllers"
	"vix-btpns/database"
	. "vix-btpns/middlewares"

	"github.com/gin-gonic/gin"
)

type Route struct{}

func (r *Route) Init() *gin.Engine {
	db := database.GetDB()

	// controllers
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	photoController := controllers.NewPhotoController(db)

	route := gin.Default()
	route.Static("/images", "./public/images")

	// API Versioning
	api := route.Group("/api/v1")

	userRoute := api.Group("/users") 
	{
		userRoute.POST("/register", authController.Register)
		userRoute.GET("/login", authController.Login)
		userRoute.PUT("/:userId", AuthMiddleware(db), userController.Update)
		userRoute.DELETE("/:userId", AuthMiddleware(db), userController.Delete)
	}

	photoRoute := api.Group("/photos") 
	{
		photoRoute.POST("", AuthMiddleware(db), photoController.Create)
		photoRoute.GET("", AuthMiddleware(db), photoController.Get)
		photoRoute.PUT("", AuthMiddleware(db), photoController.Update)
		photoRoute.DELETE("", AuthMiddleware(db), photoController.Delete)
	}

	route.Run()

	return route
}