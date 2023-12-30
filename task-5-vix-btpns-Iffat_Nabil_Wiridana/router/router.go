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
	authController := controllers.NewAuth(db)
	userController := controllers.NewUserController(db)
	photoController := controllers.NewPhoto(db)

	route := gin.Default()
	route.Static("/images", "./public/images")

	// API Versioning
	api := route.Group("/api/v1")

	userRoute := api.Group("/users")
	{
		userRoute.POST("/register", authController.SignUp)
		userRoute.GET("/login", authController.SignIn)
		userRoute.PUT("/:userId", AuthMiddleware(db), userController.Update)
		userRoute.DELETE("/:userId", AuthMiddleware(db), userController.Delete)
	}

	photoRoute := api.Group("/photos")
	{
		photoRoute.POST("", AuthMiddleware(db), photoController.Add)
		photoRoute.GET("", AuthMiddleware(db), photoController.Fetch)
		photoRoute.PUT("", AuthMiddleware(db), photoController.Modify)
		photoRoute.DELETE("", AuthMiddleware(db), photoController.Remove)
	}

	route.Run()

	return route
}
