package routes

import (
	"ghostbox/user-service/controllers"
	"github.com/buaazp/fasthttprouter"
)

func InitalizeUserRoutes(router *fasthttprouter.Router){
	router.POST("/user", controllers.CreateUser)
	router.GET("/user/:id", controllers.GetUser)
	router.PUT("/user/:id", controllers.UpdateUser)
	router.DELETE("/user/:id", controllers.DeleteUser)
}