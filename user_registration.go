package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type UserRegisterResponse struct {
	User User `json:"user"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	User UserRequest `json:"user"`
}

func Token(w http.ResponseWriter, r *http.Request) {
	req := &UserTokenRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := UserDB().GetUser(req.Email)
	if user == nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if CompareHashAndPasswordUtil([]byte(user.HashedPassword), []byte(req.Password)) {
		log.Printf("Token :: User login attempt (%s) successful!", user.Email)
		trt, err := CreateTokenResponse(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response(w, trt)
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	req := &UserRegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := UserDB().GetUser(req.User.Email)
	if user != nil {
		http.Error(w, errors.New("User already exist").Error(), http.StatusConflict)
		return
	}

	user = &User{
		Email: req.User.Email,
		ID:    int64(uuid.New().ID()),
	}
	user.HashedPassword = GenerateBCryptPasswordUtil([]byte(req.User.Password))
	UserDB().AddUser(user)

	toReturn := &UserRegisterResponse{}
	toReturn.User = User{
		Email:          user.Email,
		HashedPassword: "*********",
		ID:             user.ID,
	}

	response(w, toReturn)
}
