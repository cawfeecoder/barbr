package models

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	FirstName string `validate:"required" bson:"fname"`
	LastName string `validate:"required" bson:"lname"`
	Email string `validate:"required,email" bson:"email"`
	Password string `validate:"required,min=8" bson:"password"`
}

func (u *User) ConvertToDTO() *UserDTO{
	return &UserDTO {
		FirstName: u.FirstName,
		LastName: u.LastName,
		Email: u.Email,
	}
}

type UserDTO struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"fname" bson:"fname"`
	LastName string `json:"lname" bson:"lname"`
	Email string `json:"email" bson:"email"`
}
