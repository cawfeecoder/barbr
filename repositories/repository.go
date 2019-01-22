package repositories

import "github.com/mongodb/mongo-go-driver/mongo"

type Repository interface {
	GetCollection() (*mongo.Collection)
	Create(data interface{}) (interface{}, error)
	Get(id string) (interface{}, error)
	Update(id string, data interface{}) (interface{}, error)
	Delete(id string) (bool, error)
}

type QueryHandler func(arg []interface{}, collection *mongo.Collection) (result interface{}, err error)