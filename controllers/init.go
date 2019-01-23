package controllers

import (
	"ghostbox/user-service/repositories"
	"go.uber.org/zap"
)

var logger *zap.Logger
var user_repo repositories.UserRepository

func InitializeControllers(ur repositories.UserRepository){
	logger, _ = zap.NewProduction()
	user_repo = ur
}