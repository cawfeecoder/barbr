package routes

import (
	"ghostbox/user-service/controllers"
	"github.com/buaazp/fasthttprouter"
)

func InitializeAuthRoutes(router *fasthttprouter.Router){
	router.POST("/login", controllers.Authenticate)
}