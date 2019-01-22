package controllers

import (
	"encoding/json"
	"ghostbox/user-service/models"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func InternalServerError(ctx *fasthttp.RequestCtx){
	ctx.SetStatusCode(500)
	ctx.SetBody([]byte("Internal Server Error"))
}

func Unauthorized(ctx *fasthttp.RequestCtx, payload []byte){
	ctx.SetStatusCode(400)
	ctx.SetBody(payload)
}

func HandleUnmarshal(ctx *fasthttp.RequestCtx, target interface{}) (err error){
	err = json.Unmarshal(ctx.Request.Body(), &target)
	if err != nil {
		logger.Error("failed to unmarshal body", zap.Error(err))
		InternalServerError(ctx)
	}
	return
}

func HandleMarshal(ctx *fasthttp.RequestCtx, source interface{}) (target []byte){
	target, err := json.Marshal(source)
	if err != nil {
		InternalServerError(ctx)
		return nil
	}
	return target
}

func HandleErrors(ctx *fasthttp.RequestCtx, errors []models.HumanReadableStatus){
	resp := models.Response{}
	resp.Errors = append(resp.Errors, errors)
	json, err := resp.ToJSON()
	if err != nil {
		logger.Error("could not marshal to json", zap.Error(err))
		InternalServerError(ctx)
		return
	}
	Unauthorized(ctx, json)
	return
}