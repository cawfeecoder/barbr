package main

import (
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
	user_repository = repositories.InitalizeUserRepository()
	user_repository.EnsureIndex()
	router := fasthttprouter.New()
	routes.InitalizeUserRoutes(router, user_repository)

	logger.Info("starting server", zap.String("port", "12345"))

	logger.Fatal("error starting server", zap.Error(fasthttp.ListenAndServe(":12345", router.Handler)))
}
