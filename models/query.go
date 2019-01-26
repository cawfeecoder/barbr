package models

import (
	"github.com/modern-go/reflect2"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

type Query struct {
	Fields []string `validate:"min=0" json:"first,omitempty"`
	First int `validate:"numeric,min=0,max=50" json:"first,omitempty"`
	After string `validate:"min=0" json:"after,omitempty"`
	Cond map[string][]map[string]bson.M `json:"cond,omitempty"`
	//TODO: The above should really be a map[string]interface{} which then performs type reflection in order to determine if the
	// target is a value (equality), object (single field conditions), array(logical and/or/nor). Notice that an object can be recursive (logical -> logical -> object -> array). We will have to figure out a nicer way to do this
	// We also need to initialize the map to have the after property if it is zero, which should be easy enough.
}

func (q *Query) ConvertIDs() (errors []HumanReadableStatus) {
	if len(q.Cond) == 0 {
		q.Cond["$and"]= []map[string]bson.M {
			map[string]bson.M{}
		}
	}
	for key, val := range q.Cond {
		for k, v := range val {
			if val2, ok := v["_id"]; ok {
				for kk, vv := range val2 {
					if reflect2.TypeOf(vv).Kind() == reflect.String {
						obj_id, err := primitive.ObjectIDFromHex(vv.(string))
						if err != nil {
							errors = append(errors, HumanReadableStatus{Type: "not-valid-query", Message: "The provided valid is not a valid id", Value: vv.(string), Param: "id"})
							return
						}
						q.Cond[string(key)][k]["_id"][string(kk)] = obj_id
						if len(q.After) > 0 {
							obj_id, err := primitive.ObjectIDFromHex(q.After)
							if err != nil {
								errors = append(errors, HumanReadableStatus{Type: "not-valid-pagination", Message: "The provided after value is not a valid id", Value: q.After, Param: "after"})
								return
							}
							q.Cond[string(key)][k]["_id"]["$gt"] = obj_id
						}
					} else {
						for kkk, vvv := range vv.([]interface{}) {
							obj_id, err := primitive.ObjectIDFromHex(vvv.(string))
							if err != nil {
								errors = append(errors, HumanReadableStatus{Type: "not-valid-query", Message: "The provided valid is not a valid id", Value: vv.(string), Param: "id"})
								continue
							}
							q.Cond[string(key)][k]["_id"][string(kk)].([]interface{})[kkk] = obj_id
						}
						if len(q.After) > 0 {
							obj_id, err := primitive.ObjectIDFromHex(q.After)
							if err != nil {
								errors = append(errors, HumanReadableStatus{Type: "not-valid-pagination", Message: "The provided after value is not a valid id", Value: q.After, Param: "after"})
								return
							}
							q.Cond[string(key)][k]["_id"]["$gt"] = obj_id
						}
						return
					}
				}
			}
		}
	}
	return
}

func (q *Query) Validate(validate *validator.Validate, param string) (err []HumanReadableStatus){
	var human_readable_err ValidationErrors
	validation_err := validate.Struct(q)
	if validation_err != nil {
		human_readable_err = ValidationErrors{Err: validation_err.(validator.ValidationErrors)}
		err = human_readable_err.ToHumanReadable(param)
	}
	return
}