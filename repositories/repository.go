package repositories

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Repository interface {
	GetCollection() *mongo.Collection
	Execute(arg []interface{}, param string, q QueryHandler) (interface{}, error)
	Create(data interface{}) (interface{}, error)
	Get(id string) (interface{}, error)
	Update(id string, data interface{}) (interface{}, error)
	Delete(id string) (bool, error)
}

type QueryHandler func(arg []interface{}, projection map[string]int, collection *mongo.Collection) (result interface{}, err error)

func GenerateProjectionFromFields(fields []string) (map[string]int){
	projection := make(map[string]int)
	projection["_id"] = 0
	for _, val := range fields {
		if val == "id" {
			projection["_id"] = 1
		} else if val != "" {
			projection[val] = 1
		}
	}
	return projection
}