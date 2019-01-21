package repositories

import (
	"context"
	"fmt"
	"ghostbox/user-service/db"
	"ghostbox/user-service/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type UserRepository struct {
	database string
	collectionName string
	validator *validator.Validate
	logger *zap.Logger
}

func InitalizeUserRepository() (repository UserRepository){
	logger, _ := zap.NewProduction()
	return UserRepository{
		database: "barbr",
		collectionName: "users",
		validator: validator.New(),
		logger: logger,
	}
}

func (r *UserRepository) GetUserCollection() *mongo.Collection{
	return db.GetClient().Database(r.database).Collection(r.collectionName)
}

func (r *UserRepository) Create(user *models.User) (result *models.UserDTO, err error) {
	c := r.GetUserCollection()
	err = r.validator.Struct(user)
	if err != nil {
		r.logger.Error("cannot validate new user", zap.Error(err))
		return nil, err
	}
	res, err := c.InsertOne(context.Background(), user)
	if err != nil {
		r.logger.Error("cannot insert new user", zap.Error(err))
		return nil, err
	}
	result = user.ConvertToDTO()
	result.ID = res.InsertedID.(primitive.ObjectID)
	return result, nil
}

func (r *UserRepository) Get(id string) (result *models.UserDTO, err error){
	c := r.GetUserCollection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return nil, err
	}
	user := models.UserDTO{}
	err = c.FindOne(context.Background(), bson.D{{"_id",object_id}}).Decode(&user)
	if err != nil {
		fmt.Printf(err.Error())
		r.logger.Error("failed to decode data", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
