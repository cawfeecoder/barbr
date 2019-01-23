package main

import (
	"ghostbox/user-service/controllers"
	"ghostbox/user-service/repositories"
	"ghostbox/user-service/routes"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var logger *zap.Logger

var user_repository repositories.UserRepository

func main() {
	logger, _ = zap.NewProduction()
	logger.Info("initializing repository", zap.String("repository", "user"))
	user_repository = repositories.InitalizeUserRepository()
	logger.Info("initializing index", zap.String("index", "user"))
	user_repository.EnsureIndex()
	logger.Info("initializing fasthttp router")
	router := fasthttprouter.New()
	logger.Info("initializing controllers")
	controllers.InitializeControllers(user_repository)
	logger.Info("initializing routers", zap.String("routes", "user"))
	routes.InitalizeUserRoutes(router, user_repository)
	logger.Info("starting server", zap.String("port", "12345"))
	logger.Fatal("error starting server", zap.Error(fasthttp.ListenAndServe(":12345", router.Handler)))
}
