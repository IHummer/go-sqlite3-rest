package main

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ApiError struct {
	Error string `json:"error"`
}

type User struct {
	Id             int
	Username       string
	Email          string
	Password       string
	HashedPassword string
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newAccount(firstName, lastName string) *Account {
	return &Account{
		// Id:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    uuid.New().String(),
		CreatedAt: time.Now().UTC(),
	}
}
