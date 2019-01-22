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
	r.GetCollection().Indexes().CreateOne(context.Background(), index, opts)
	r.logger.Info("Successfully created index", zap.String("index", "email-uniq"))
}

func (r *UserRepository) GetCollection() *mongo.Collection{
	return db.GetClient().Database(r.database).Collection(r.collectionName)
}

func (r *UserRepository) Authenticate(email string, password []byte) (result bool, err error) {
	c := r.GetCollection()
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

func (r *UserRepository) Execute(arg []interface{}, param string, q QueryHandler) (result interface{}, errs []models.HumanReadableStatus) {
	c := r.GetCollection()
	result, err := q(arg, c)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err))
		return nil, models.GetErrorFromMongo(err, param)
	}
	return
}

func (r *UserRepository) Create(arg []interface{}, c *mongo.Collection) (result interface{}, err error) {
	user := arg[0].(models.User)
	user.New()
	res, err := c.InsertOne(context.Background(), user)
	if err != nil {
		r.logger.Error("failed to insert data", zap.Error(err))
		return
	}
	result = user.ConvertToDTO(res.InsertedID)
	return
}

func (r *UserRepository) Get(arg []interface{}, c *mongo.Collection) (result interface{}, err error){
	var user models.UserDTO
	var object_id primitive.ObjectID
	if object_id, err = primitive.ObjectIDFromHex(arg[0].(string)); err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return
	}
	if err = c.FindOne(context.Background(), bson.D{{"_id",object_id}}).Decode(&user); err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return
	}
	result = user
	return
}

func (r *UserRepository) Update(arg []interface{}, c *mongo.Collection) (result interface{}, err error){
	var user models.UserDTO
	object_id, err := primitive.ObjectIDFromHex(arg[0].(string))
	if err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return
	}
	opts := options.FindOneAndUpdateOptions{}
	opts.SetReturnDocument(options.After)
	err = c.FindOneAndUpdate(context.Background(), bson.D{{"_id", object_id}}, bson.D{{"$set", arg[1].(models.UserDTO)}}, &opts).Decode(&user)
	if err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return
	}
	result = user
	return
}

func (r *UserRepository) Delete(arg []interface{}, c *mongo.Collection) (result interface{}, err error){
	var user models.UserDTO
	object_id, err := primitive.ObjectIDFromHex(arg[0].(string))
	if err != nil {
		r.logger.Error("cannot convert to object id", zap.Error(err))
		return
	}
	opts := options.FindOneAndUpdateOptions{}
	opts.SetReturnDocument(options.After)
	err = c.FindOneAndUpdate(context.Background(), bson.D{{"_id", object_id}}, bson.D{{"$set", bson.M{"status": "purge"}}}).Decode(&user)
	if err != nil {
		r.logger.Error("failed to decode data", zap.Error(err))
		return
	}
	result = user
	return
}

