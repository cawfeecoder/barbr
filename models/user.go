package models

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	Status string `validate:"isdefault" json:"status" bson:"status"`
	FName string `validate:"required" json:"fname" bson:"fname"`
	LName string `validate:"required" json:"lname" bson:"lname"`
	Email string `validate:"required,email" json:"email" bson:"email"`
	Password string `validate:"required,min=8" json:"password" bson:"password"`
}

func (u *User) HashPassword() {
	cfg := argon2.DefaultConfig()
	bytes, err := cfg.HashEncoded([]byte(u.Password))
	if err != nil {
		fmt.Printf("failed to hash password: %v", err)
	}
	u.Password = string(bytes)
}

func (u *User) ConvertToDTO() *UserDTO{
	return &UserDTO {
		Status: u.Status,
		FName: u.FName,
		LName: u.LName,
		Email: u.Email,
	}
}

type UserDTO struct {
	ID primitive.ObjectID `validate:"isdefault" json:"id,omitempty" bson:"_id,omitempty"`
	Status string `validate:"isdefault" json:"status,omitempty" bson:"status,omitempty"`
	FName string `json:"fname,omitempty" bson:"fname,omitempty"`
	LName string `json:"lname,omitempty" bson:"lname,omitempty"`
	Email string `validate:"email" json:"email,omitempty" bson:"email,omitempty"`
}
