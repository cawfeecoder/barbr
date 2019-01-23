package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"ghostbox/user-service/models"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func ExtractAuthenticationHeader(ctx *fasthttp.RequestCtx) (email string, password []byte, auth_err []models.HumanReadableStatus){
	auth := ctx.Request.Header.Peek("Authorization")
	if bytes.HasPrefix(auth, []byte("Basic ")) {
		payload, err := base64.StdEncoding.DecodeString(string(auth[len([]byte("Basic ")):]))
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 {
				email = string(pair[0])
				password = pair[1]
			}
		} else {
			logger.Error("failed to marshal data", zap.Error(err))
			InternalServerError(ctx)
			return
		}
	} else {
		err := errors.New("Basic is not present in authorization header")
		logger.Error("missing authorization header", zap.Error(err))
		auth_err = append(auth_err, models.HumanReadableStatus{Type: "user-no-credentials", Message: "No credentials provided for user", Param: "username, password"})
	}
	return
}

func Authenticate(ctx *fasthttp.RequestCtx) {
	projection := make(map[string]int)
	resp := models.Response{}
	email, password, err := ExtractAuthenticationHeader(ctx)
	if err != nil {
		logger.Error("failed to authenticate", zap.Errors("errors", models.ToErrorsArray(err)))
		ctx.SetStatusCode(401)
		HandleErrors(ctx, resp, err)
		return
	}
	res, err := user_repo.Execute([]interface{}{email, password}, "password", projection, user_repo.Authenticate)
	if err != nil || !res.(bool) {
		logger.Error("failed to authenticate", zap.Errors("errors", models.ToErrorsArray(err)))
		HandleErrors(ctx, resp, err)
		return
	}
	resp.Data = []interface{}{models.HumanReadableStatus{Type:"user-auth-success", Message: fmt.Sprintf("User %s has successfully authenticated", string(email))}}
	if data := HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}