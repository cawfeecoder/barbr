package routes

import (
	"ghostbox/user-service/controllers"
	"ghostbox/user-service/models"
	"ghostbox/user-service/repositories"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var user_repo repositories.UserRepository
var logger *zap.Logger
var validate *validator.Validate

func InitalizeUserRoutes(router *fasthttprouter.Router, user_repository repositories.UserRepository){
	user_repo = user_repository
	logger, _ = zap.NewProduction()
	validate = validator.New()
	router.POST("/login", controllers.Authenticate)
	router.POST("/user", CreateUser)
	router.GET("/user/:id", GetUser)
	router.PUT("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)
}

func CreateUser(ctx *fasthttp.RequestCtx) {
	projection := make(map[string]int)
	if ctx.QueryArgs().Has("fields") {
		params := strings.Split(string(ctx.QueryArgs().Peek("fields")), ",")
		projection = repositories.GenerateProjectionFromFields(params)
	}
	var user models.User
	resp := models.Response{}
	err := controllers.HandleUnmarshal(ctx, &user)
	if err != nil {
		return
	}
	validation_errs := user.Validate(validate,"")
	if len(validation_errs) > 0 {
		controllers.HandleErrors(ctx, resp, validation_errs)
		return
	}
	res, query_err := user_repo.Execute([]interface{}{user}, "", projection, user_repo.Create)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, resp, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func GetUser(ctx *fasthttp.RequestCtx) {
	projection := make(map[string]int)
	if ctx.QueryArgs().Has("fields") {
		params := strings.Split(string(ctx.QueryArgs().Peek("fields")), ",")
		projection = repositories.GenerateProjectionFromFields(params)
	}
	resp := models.Response{}
	res, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id")}, ctx.UserValue("id").(string), projection, user_repo.Get)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, resp, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	projection := make(map[string]int)
	if ctx.QueryArgs().Has("fields") {
		params := strings.Split(string(ctx.QueryArgs().Peek("fields")), ",")
		projection = repositories.GenerateProjectionFromFields(params)
	}
	var user models.UserDTO
	resp := models.Response{}
	err := controllers.HandleUnmarshal(ctx, &user)
	if err != nil {
		return
	}
	validation_errs := user.Validate(validate,"")
	if len(validation_errs) > 0 {
		controllers.HandleErrors(ctx, resp, validation_errs)
		return
	}
	res, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id"), user}, ctx.UserValue("id").(string), projection, user_repo.Update)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, resp, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func DeleteUser(ctx *fasthttp.RequestCtx) {
	projection := make(map[string]int)
	resp := models.Response{}
	_, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id").(string)}, "id", projection, user_repo.Delete)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, resp, query_err)
		return
	}
	resp.Data = []interface{}{models.HumanReadableStatus{Type: "account-delete-success", Message: "Your account has successfully been marked for deletion and will be purged within 72 hours.", Source: ctx.UserValue("id").(string)}}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}