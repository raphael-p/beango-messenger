package mocks

import (
	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

func MakeUser() *database.User {
	return &database.User{
		ID:          uuid.New().String(),
		Username:    "john doe",
		DisplayName: "john.doe.69",
		Key:         []byte("supersecrethash"),
	}
}
