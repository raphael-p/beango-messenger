package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

const ADMIN_USERNAME = "the_admin"
const PASSWORD = "123abc*"
const HASH = "$2y$04$8QoTLjUMGtnr4lNeA0DtduhEshmvbDbmEzW/G9IkV/9mr576xX//K"

func MakeUser() *database.User {
	return &database.User{
		ID:          uuid.New().String(),
		Username:    "john.doe.69",
		DisplayName: "Johnny D",
		Key:         []byte(HASH),
	}
}

func MakeAdminUser() *database.User {
	return &database.User{
		ID:          uuid.New().String(),
		Username:    ADMIN_USERNAME,
		DisplayName: "Administrator",
		Key:         []byte(HASH),
	}
}

func MakeSession(userId string) database.Session {
	return database.Session{
		ID:         uuid.New().String(),
		UserID:     userId,
		ExpiryDate: time.Now().UTC().Add(time.Hour),
	}
}
