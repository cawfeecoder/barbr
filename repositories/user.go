package repositories

import (
	"context"
	"errors"
	"ghostbox/user-service/db"
	"ghostbox/user-service/models"
	"github.com/matthewhartstonge/argon2"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"time"
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

func (r *UserRepository) EnsureIndex() {
	index := mongo.IndexModel{}
	index.Keys = bson.M{"email": 1}
	unique := true
	index.Options = &options.IndexOptions {
		Unique: &unique,
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	r.GetUserCollection().Indexes().CreateOne(context.Background(), index, opts)
	r.logger.Info("Successfully created index", zap.String("index", "email-uniq"))
}

func (r *UserRepository) GetUserCollection() *mongo.Collection{
	return db.GetClient().Database(r.database).Collection(r.collectionName)
}

func (r *UserRepository) Authenticate(email string, password []byte) (result bool, err error) {
	c := r.GetUserCollection()
	var user models.User
	err = c.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return false, err
	}
	ok, err := argon2.VerifyEncoded(password, []byte(user.Password))
	if ok {
		return true, nil
	} else {
		return false, errors.New("Incorrect Password")
	}
}

func (r *UserRepository) Create(user *models.User) (result *models.UserDTO, err error) {
	c := r.GetUserCollection()
	user.Status = "active"
	user.HashPassword()
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
		r.logger.Error("failed to decode data", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(id string, update models.UserDTO) (result *models.UserDTO, err error){
	c := r.GetUserCollection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return nil, err
	}
	return_opts := options.FindOneAndUpdateOptions{}
	return_opts.SetReturnDocument(options.After)
	err = c.FindOneAndUpdate(context.Background(), bson.D{{"_id", object_id}}, bson.D{{"$set", update}}, &return_opts).Decode(&result)
	if err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) Delete(id string) (result *models.UserDTO, err error){
	c := r.GetUserCollection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return nil, err
	}
	return_opts := options.FindOneAndUpdateOptions{}
	return_opts.SetReturnDocument(options.After)
	err = c.FindOneAndUpdate(context.Background(), bson.D{{"_id", object_id}}, bson.D{{"$set", bson.M{"status": "purge"}}}).Decode(&result)
	if err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return nil, err
	}
	return result, nil
}

