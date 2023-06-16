package main

import (
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	uid := GetUserId(r)

	if uid == 0 {
		http.Error(w, "user ID can't be 0", http.StatusBadRequest)
		return
	}

	user := UserDB().GetUserById(uid)
	if user == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	toReturn := &UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	response(w, toReturn)
}
