package db

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"time"
)

var client *mongo.Client
var logger *zap.Logger

func init() {
	var err error
	logger, _ := zap.NewProduction()
	client, err = mongo.NewClientWithOptions("mongodb://root:example@localhost:27017", )
	if err != nil {
		logger.Error("failed to create client", zap.Error(err))
	}
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Error("failed to connect", zap.Error(err))
	}
}

func GetClient() *mongo.Client {
	return client
}