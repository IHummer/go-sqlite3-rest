package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type SQLiteStore struct {
	db *sql.DB
}

func newSqliteStore() (*SQLiteStore, error) {
	dbfileLocation := "local.db"

	file, err := os.OpenFile(dbfileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	db, err := sql.Open("sqlite3", dbfileLocation)
	if err != nil {
		fmt.Println("test")
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Init() error {
	return s.CreateAccountTable()
}
func (s *SQLiteStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number VARCHAR(300),
		balance DECIMAL(10,2),
		created_at DATETIME
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) CreateAccount(account *Account) (int, error) {
	// return nil
	query := `
	INSERT INTO account(first_name, last_name, number, balance, created_at) 
	VALUES (?, ?, ?, ?, ?)`
	resp, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SQLiteStore) UpdateAccount(*Account) error {
	return nil
}
func (s *SQLiteStore) DeleteAccount(id int) error {
	return nil
}
func (s *SQLiteStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
