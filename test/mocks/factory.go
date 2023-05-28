package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

const ADMIN_ID = "09d3fb32-af60-4029-b37e-a3b92dc0a559"
const ADMIN_USERNAME = "the_admin"
const PASSWORD = "123abc*"
const HASH = "$2y$04$8QoTLjUMGtnr4lNeA0DtduhEshmvbDbmEzW/G9IkV/9mr576xX//K"

func MakeUser() *database.User {
	return &database.User{
		ID:          uuid.NewString(),
		Username:    "john.doe.69",
		DisplayName: "Johnny D",
		Key:         []byte(HASH),
	}
}

func MakeAdminUser() *database.User {
	return &database.User{
		ID:          ADMIN_ID,
		Username:    ADMIN_USERNAME,
		DisplayName: "Administrator",
		Key:         []byte(HASH),
	}
}

func MakeSession(userID string) database.Session {
	return database.Session{
		ID:         uuid.NewString(),
		UserID:     userID,
		ExpiryDate: time.Now().UTC().Add(time.Hour),
	}
}

func MakeChat(userIDOne, userIDTwo string) *database.Chat {
	return &database.Chat{
		ID:      uuid.NewString(),
		UserIDs: [2]string{userIDOne, userIDTwo},
	}
}

func MakeMessage(userID, chatID string) *database.Message {
	return &database.Message{
		ID:      uuid.NewString(),
		UserID:  userID,
		ChatID:  chatID,
		Content: "Lorem Ipsum Dolor",
	}
}
