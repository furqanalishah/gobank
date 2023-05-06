package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
    		id serial primary key,
        	first_name varchar(50),
        	last_name varchar(50),
    		number serial,
    		balance serial,
    		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(a *Account) error {
	query := `INSERT INTO accounts(first_name, last_name, number, balance, created_at)
			 VALUES ($1, $2, $3, $4, $5)`
	resp, err := s.db.Exec(query, a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id=$1", id)
	return err
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	return account, err
}
