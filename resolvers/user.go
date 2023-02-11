package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	_, vals := utils.MapValues(database.Users)
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(vals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	newUser := database.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	if input.Username == "" {
		http.Error(w, "username is missing", http.StatusBadRequest)
		return
	}

	for _, value := range database.Users {
		if value.Username == input.Username {
			http.Error(w, "username is taken", http.StatusConflict)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUser.Id = uuid.New().String()
	newUser.Username = input.Username
	newUser.Key = string(hash)
	database.Users[newUser.Id] = newUser
	w.WriteHeader(http.StatusCreated)
	response, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
