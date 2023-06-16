package main

type User struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	ID             int64  `json:"id"`
}
