package models

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	Status string `validate:"isdefault" json:"status" bson:"status"`
	FName string `validate:"required" json:"fname" bson:"fname"`
	LName string `validate:"required" json:"lname" bson:"lname"`
	Email string `validate:"required,email" json:"email" bson:"email"`
	Password string `validate:"required,min=8" json:"password" bson:"password"`
}

func (u *User) New() {
	u.Status = "active"
	u.HashPassword()
}

func (u *User) HashPassword() {
	cfg := argon2.DefaultConfig()
	bytes, err := cfg.HashEncoded([]byte(u.Password))
	if err != nil {
		fmt.Printf("failed to hash password: %v", err)
	}
	u.Password = string(bytes)
}

func (u *User) ConvertToDTO(id interface{}) *UserDTO{
	if id != nil {
		return &UserDTO {
			ID: id.(primitive.ObjectID),
			Status: u.Status,
			FName: u.FName,
			LName: u.LName,
			Email: u.Email,
		}
	}
	return &UserDTO {
		Status: u.Status,
		FName: u.FName,
		LName: u.LName,
		Email: u.Email,
	}
}

func (u *User) Validate(param string) (err []HumanReadableStatus) {
	validate := validator.New()
	var human_readable_err ValidationErrors
	validation_err := validate.Struct(u)
	if validation_err != nil {
		human_readable_err = ValidationErrors{Err: validation_err.(validator.ValidationErrors)}
		err = human_readable_err.ToHumanReadable(param)
	}
	return
}

type UserDTO struct {
	ID primitive.ObjectID `validate:"isdefault" json:"id,omitempty" bson:"_id,omitempty"`
	Status string `validate:"isdefault" json:"status,omitempty" bson:"status,omitempty"`
	FName string `json:"fname,omitempty" bson:"fname,omitempty"`
	LName string `json:"lname,omitempty" bson:"lname,omitempty"`
	Email string `validate:"email" json:"email,omitempty" bson:"email,omitempty"`
}

func (u *UserDTO) Validate(param string) (err []HumanReadableStatus) {
	validate := validator.New()
	var human_readable_err ValidationErrors
	validation_err := validate.Struct(u)
	if validation_err != nil {
		human_readable_err = ValidationErrors{Err: validation_err.(validator.ValidationErrors)}
		err = human_readable_err.ToHumanReadable(param)
	}
	return
}