package routes

import (
	"encoding/json"
	"ghostbox/user-service/models"
	"ghostbox/user-service/repositories"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var user_repo repositories.UserRepository
var logger *zap.Logger

func InitalizeUserRoutes(router *fasthttprouter.Router, user_repository repositories.UserRepository){
	user_repo = user_repository
	logger, _ = zap.NewProduction()
	router.GET("/user/:id", GetUser)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	res, err := user_repo.Get(ctx.UserValue("id").(string))
	resp := models.Response{}
	if err != nil {
		logger.Error("failed to fetch user", zap.Error(err))
		resp.Errors = []interface{}{models.GetErrorFromMongo(err)}
		resp_json, err := json.Marshal(resp)
		if err != nil {
			logger.Error("failed to marshal to []byte", zap.Error(err))
			ctx.SetStatusCode(500)
			ctx.SetBody([]byte("Internal Server Error"))
			return
		}
		ctx.SetStatusCode(404)
		ctx.SetBody([]byte(resp_json))
		return
	}
	resp.Data = []interface{}{res}
	resp_json, err := json.Marshal(resp)
	if err != nil {
		logger.Error("failed to marshal to []byte", zap.Error(err))
		ctx.SetStatusCode(500)
		ctx.SetBody([]byte("Internal Server Error"))
		return
	}
	ctx.Success("application/json", resp_json)
	return
}
