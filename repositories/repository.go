package repositories

type Repository interface {
	Create(data interface{}) (interface{}, error)
	Get(id string) (interface{}, error)
	Update(id string, data interface{}) (interface{}, error)
	Delete(id string) (bool, error)
}
