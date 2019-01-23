package routes

import "github.com/buaazp/fasthttprouter"

func InitializeRoutes(router *fasthttprouter.Router){
	InitializeAuthRoutes(router)
	InitalizeUserRoutes(router)
}