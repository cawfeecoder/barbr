package main

import (
	"encoding/json"
	"fmt"
	"ghostbox/user-service/models"
	"ghostbox/user-service/repositories"
	"ghostbox/user-service/routes"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

var logger *zap.Logger

var user_repository repositories.UserRepository

func Index(ctx *fasthttp.RequestCtx) {
	user := &models.User{
		FirstName: "John",
		LastName: "Doe",
		Email: "john.doe@gmail.com",
		Password: "abcd1234",
	}
	res, err := user_repository.Create(user)
	resp := models.Response{}
	if err != nil {
		validation_err := models.ValidationErrors{err.(validator.ValidationErrors)}
		resp.Errors = append(resp.Errors, validation_err.ToHumanReadable())
		if err != nil {
			logger.Error("could not marshal to json", zap.Error(err))
			ctx.Error(fmt.Sprintf("Error: %v", err.Error()), 500)
		}
		resp_json, err := resp.ToJSON()
		if err != nil {
			logger.Error("could not marshal to json", zap.Error(err))
			ctx.Error(fmt.Sprintf("Error: %v", err.Error()), 500)
		}
		ctx.SetStatusCode(400)
		ctx.SetBody(resp_json)
		return
	}
	resp.Data = []interface{}{res}
	resp_json, err := json.Marshal(resp)
	if err != nil {
		logger.Error("could not marshal to json", zap.Error(err))
		ctx.Error(fmt.Sprintf("Error: %v", err.Error()), 500)
	}
	ctx.Success("Success", resp_json)
}

func main() {
	logger, _ = zap.NewProduction()
	user_repository = repositories.InitalizeUserRepository()
	router := fasthttprouter.New()
	router.GET("/", Index)
	routes.InitalizeUserRoutes(router, user_repository)

	logger.Info("starting server", zap.String("port", "12345"))

	logger.Fatal("error starting server", zap.Error(fasthttp.ListenAndServe(":12345", router.Handler)))
}
