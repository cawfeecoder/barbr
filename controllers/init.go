package controllers

import (
	"ghostbox/user-service/repositories"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

var logger *zap.Logger
var user_repo repositories.UserRepository
var validate *validator.Validate

func InitializeControllers(ur repositories.UserRepository){
	logger, _ = zap.NewProduction()
	user_repo = ur
	validate = validator.New()
}