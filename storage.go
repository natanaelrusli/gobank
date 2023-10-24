package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
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

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln(err)
	}
}

func NewPostgresStore() (*PostgresStore, error) {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbUser, dbName, dbPassword)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account 
	(first_name, last_name, number, balance, created_at)
	values 
	($1, $2, $3, $4, $5)
	`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	q := "DELETE FROM account WHERE id = $1;"

	_, err := s.db.Query(q, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	// UPDATE table_name
	// SET column1 = value1, column2 = value2, ...
	// WHERE condition;
	q := "UPDATE account SET first_name = $1, last_name = $2 WHERE id = $3"

	row := s.db.QueryRow(q, account.FirstName, account.LastName, account.ID)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with id %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}
