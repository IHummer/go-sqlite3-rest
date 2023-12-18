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

func newAccount(firstName, lastName string) *Account {
	return &Account{
		// Id:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    uuid.New().String(),
		CreatedAt: time.Now().UTC(),
	}
}
