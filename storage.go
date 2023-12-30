package main

import (
	"database/sql"
	"errors"
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
	GetAllAccounts() ([]*Account, error)
	RegisterUser(*User) (int, error)
	SelectUserPassword(string) (string, error)
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
	return s.CreateInitTables()
}
func (s *SQLiteStore) CreateInitTables() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number VARCHAR(300),
		balance DECIMAL(10,2),
		created_at DATETIME
	)
	;
	CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE,
		email VARCHAR(250),
		password VARCHAR(300),
		created_at DATETIME
	)
	;
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) RegisterUser(user *User) (int, error) {
	// return nil
	query := `
	INSERT INTO User(username, email, password, created_at) 
	VALUES (?, ?, ?, ?)`
	resp, err := s.db.Exec(query, user.Username, user.Email, user.HashedPassword, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SQLiteStore) SelectUserPassword(username string) (string, error) {
	// return nil
	password := ""
	query := "SELECT password FROM User WHERE username = ?"
	resp := s.db.QueryRow(query, username)

	err := resp.Scan(&password)
	if err != nil && err == sql.ErrNoRows {
		return password, nil
	} else if err != nil && err != sql.ErrNoRows {
		return password, err
	}
	return password, nil
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

func (s *SQLiteStore) UpdateAccount(account *Account) error {
	exists, err := s.AccountExists(account.Id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Objeto no encontrado")
	}

	query := `
	UPDATE account
	SET first_name = ?, 
		last_name = ?, 
		number = ?, 
		balance = ? 
	WHERE id = ?`

	_, err = s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.Id)

	return err
}
func (s *SQLiteStore) DeleteAccount(id int) error {
	exists, err := s.AccountExists(id)
	if err != nil {
		//db error
		return err
	}
	if !exists {
		return errors.New("Objeto no encontrado")
	}
	query := "DELETE FROM account WHERE id = ?"
	_, err = s.db.Exec(query, id)
	return err
}

func (s *SQLiteStore) GetAllAccounts() ([]*Account, error) {
	query := `SELECT 
		id AS Id,
		first_name AS FirstName, 
		last_name AS LastName, 
		number AS Number, 
		balance AS Balance, 
		created_at AS CreatedAt
	FROM account`
	// account := new(Account)
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.Id,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *SQLiteStore) GetAccountById(id int) (*Account, error) {
	query := `SELECT 
		id AS Id,
		first_name AS FirstName, 
		last_name AS LastName, 
		number AS Number, 
		balance AS Balance, 
		created_at AS CreatedAt
	FROM account WHERE id = ?`
	// account := new(Account)
	row := s.db.QueryRow(query, id)

	account, err := ScanIntoAccount(row)

	if err != nil && err != sql.ErrNoRows {
		//db error
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		//account doesn't exists
		return account, nil
	}

	return account, nil
}

func (s *SQLiteStore) AccountExists(id int) (bool, error) {
	idScan := 0

	query := "SELECT id FROM account WHERE id = ?"
	row := s.db.QueryRow(query, id)
	err := row.Scan(&idScan)

	if err != nil && err != sql.ErrNoRows {
		//db error
		return false, err
	} else if err != nil && err == sql.ErrNoRows {
		//account doesn't exists
		return false, nil
	}
	return true, nil
}

func ScanIntoAccount(rows *sql.Row) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.Id,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	return account, err
}
