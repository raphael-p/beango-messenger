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
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vals)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	newUser := database.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "malformed request body")
		return
	}

	if input.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "username is missing")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	newUser.Id = uuid.New().String()
	newUser.Username = input.Username
	newUser.Key = string(hash)
	database.Users[newUser.Id] = newUser
	w.WriteHeader(http.StatusCreated)
	response, err := json.Marshal(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "could not encode response object")
	}
	w.Write(response)
}
