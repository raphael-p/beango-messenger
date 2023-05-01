package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

func MakeUser() *database.User {
	return &database.User{
		ID:          uuid.New().String(),
		Username:    "john.doe.69",
		DisplayName: "Johnny D",
		Key:         []byte("supersecrethash"),
	}
}

func MakeAdminUser() *database.User {
	return &database.User{
		ID:          uuid.New().String(),
		Username:    "the_admin",
		DisplayName: "Administrator",
		Key:         []byte("123abc*"),
	}
}

func MakeSession(userId string) database.Session {
	return database.Session{
		ID:         uuid.New().String(),
		UserID:     userId,
		ExpiryDate: time.Now().UTC().Add(time.Hour),
	}
}
