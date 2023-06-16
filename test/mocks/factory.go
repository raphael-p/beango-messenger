package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

const (
	ADMIN_ID       = 1
	ADMIN_USERNAME = "the_admin"
	PASSWORD       = "123abc*"
	HASH           = "$2y$04$8QoTLjUMGtnr4lNeA0DtduhEshmvbDbmEzW/G9IkV/9mr576xX//K"
)

func MakeUser(id int) *database.User {
	return &database.User{
		ID:          id,
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

func MakeSession(userID int) database.Session {
	return database.Session{
		ID:         uuid.NewString(),
		UserID:     userID,
		ExpiryDate: time.Now().UTC().Add(time.Hour),
	}
}

func MakePrivateChat(id int) *database.Chat {
	return &database.Chat{
		ID:       id,
		ChatType: database.PRIVATE_CHAT,
	}
}

func MakeMessage(messageID, userID, chatID int) *database.Message {
	return &database.Message{
		ID:      messageID,
		UserID:  userID,
		ChatID:  chatID,
		Content: "Lorem Ipsum Dolor",
	}
}
