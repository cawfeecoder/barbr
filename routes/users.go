package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"ghostbox/user-service/controllers"
	"ghostbox/user-service/models"
	"ghostbox/user-service/repositories"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

var user_repo repositories.UserRepository
var logger *zap.Logger
var validate *validator.Validate

func InitalizeUserRoutes(router *fasthttprouter.Router, user_repository repositories.UserRepository){
	user_repo = user_repository
	logger, _ = zap.NewProduction()
	validate = validator.New()
	router.POST("/login", Login)
	router.POST("/user", CreateUser)
	router.GET("/user/:id", GetUser)
	router.PUT("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)
}

/**
 * Authenticates a user by checking the hashed password against the one stored on the database record. Will return an error
 * if no user with the email exists or if there are no active user records with the email. The credentials must be provided
 * via the Authorization header to maintain best practice
 * @param ctx - HTTP Context
 * @return Unauthorized if process fails, potentially with an error. Otherwise returns true with a relevant status
 **/
func Login(ctx *fasthttp.RequestCtx) {
	auth := ctx.Request.Header.Peek("Authorization")
	var email []byte
	resp := models.Response{}
	if bytes.HasPrefix(auth, []byte("Basic ")) {
		payload, err := base64.StdEncoding.DecodeString(string(auth[len([]byte("Basic ")):]))
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 {
				email = pair[0]
				password := pair[1]
				res, err := user_repo.Authenticate(string(email), password)
				if err != nil || !res {
					logger.Error("failed to authenticate", zap.Error(err))
					ctx.SetStatusCode(401)
					resp.Errors = []interface{}{models.HumanReadableStatus{Type:"user-incorrect-password", Message:"Incorrect password was provided", Param:"password"}}
					bytes, err := json.Marshal(resp)
					if err != nil {
						logger.Error("failed to marshal body", zap.Error(err))
						ctx.SetStatusCode(500)
						ctx.SetBody([]byte("Internal Server Error"))
						return
					}
					ctx.SetBody(bytes)
					return
				}
				ctx.SetStatusCode(200)
				resp.Data = []interface{}{models.HumanReadableStatus{Type:"user-auth-success", Message: fmt.Sprintf("User %s has successfully authenticated", string(email))}}
				bytes, err := json.Marshal(resp)
				if err != nil {
					logger.Error("failed to marshal body", zap.Error(err))
					ctx.SetStatusCode(500)
					ctx.SetBody([]byte("Internal Server Error"))
					return
				}
				ctx.SetBody(bytes)
				return
			}
		} else {
			logger.Error("failed to get authorization", zap.Error(err))
		}
	} else {
		logger.Error("missing authorization header", zap.Error(errors.New("Basic is not present in authorization header")))
	}
}

func CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	resp := models.Response{}
	err := controllers.HandleUnmarshal(ctx, &user)
	if err != nil {
		return
	}
	validation_errs := user.Validate("")
	if len(validation_errs) > 0 {
		controllers.HandleErrors(ctx, validation_errs)
		return
	}
	res, query_err := user_repo.Execute([]interface{}{user}, "", user_repo.Create)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func GetUser(ctx *fasthttp.RequestCtx) {
	resp := models.Response{}
	res, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id")}, ctx.UserValue("id").(string), user_repo.Get)
	if len(query_err) > 0 {
		controllers.HandleErrors(ctx, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	var user models.UserDTO
	resp := models.Response{}
	err := controllers.HandleUnmarshal(ctx, &user)
	if err != nil {
		return
	}
	validation_errs := user.Validate("")
	if len(validation_errs) > 0 {
		controllers.HandleErrors(ctx, validation_errs)
		return
	}
	res, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id"), user}, ctx.UserValue("id").(string), user_repo.Update)
	if len(query_err) > 0 {
		fmt.Printf("Query Errors: %v", query_err)
		controllers.HandleErrors(ctx, query_err)
		return
	}
	resp.Data = []interface{}{res}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}

func DeleteUser(ctx *fasthttp.RequestCtx) {
	resp := models.Response{}
	_, query_err := user_repo.Execute([]interface{}{ctx.UserValue("id").(string)}, "id", user_repo.Delete)
	if len(query_err) > 0 {
		fmt.Printf("Query Errors: %v", query_err)
		controllers.HandleErrors(ctx, query_err)
		return
	}
	resp.Data = []interface{}{models.HumanReadableStatus{Type: "account-delete-success", Message: "Your account has successfully been marked for deletion and will be purged within 72 hours.", Source: ctx.UserValue("id").(string)}}
	if data := controllers.HandleMarshal(ctx, resp); data != nil {
		ctx.Success("Success", data)
	}
	return
}