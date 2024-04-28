package datastore

import (
	"database/sql"
	"fmt"
	"github.com/GetterSethya/golangApiMarketplace/config"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() *PostgresStorage {
	cfg := config.LoadConfig().Postgres
	connString := createPgConnStr(cfg)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Cannot establish database connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed when pinging postgres database:", err)
	}

	log.Println("Database connected successfully")

	return &PostgresStorage{
		db: db,
	}

}

func (s *PostgresStorage) Init() (*sql.DB, error) {

	// bikin tabel user
	if err := s.createUserTable(); err != nil {
		return nil, err
	}

	// bikin tabel product
	if err := s.createProductTable(); err != nil {
		return nil, err
	}

	// bikin tabel bankAccount
	if err := s.createBankAccountTable(); err != nil {
		return nil, err
	}

	// bikin tabel transaction
	if err := s.createTransactionTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *PostgresStorage) createTransactionTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions (
            id uuid NOT NULL PRIMARY KEY,
            status VARCHAR(25) NOT NULL,
            productId uuid NOT NULL,
            buyerId uuid NOT NULL,
            sellerId uuid NOT NULL,
            quantity INTEGER NOT NULL,
            notes TEXT,
            total NUMERIC(100,2) NOT NULL,
            
            createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            deletedAt TIMESTAMP
        )`)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) createBankAccountTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS bankAccounts (
            id uuid NOT NULL PRIMARY KEY,
            bankName VARCHAR(50) NOT NULL,
            accountName VARCHAR(100) NOT NULL,
            accountNumber BIGINT NOT NULL,
            sellerId uuid NOT NULL,
        
            createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            deletedAt TIMESTAMP
        )`)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) createProductTable() error {

	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id uuid NOT NULL PRIMARY KEY,
            name VARCHAR(200) NOT NULL,
            price NUMERIC(100,2) NOT NULL,
            imageUrl VARCHAR(255) NOT NULL,
            condition VARCHAR(5) NOT NULL,
            tags VARCHAR(50)[],
            isPurchaseable BOOLEAN NOT NULL DEFAULT TRUE,
            sellerId uuid NOT NULL,
            stock SMALLINT NOT NULL,
            descriptions TEXT DEFAULT '',

            createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            deletedAt TIMESTAMP
        )`)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) createUserTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id uuid NOT NULL PRIMARY KEY,
            name VARCHAR(50) NOT NULL,
            username VARCHAR(15) NOT NULL UNIQUE,
            hashPassword VARCHAR(255) NOT NULL,

            createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            deletedAt TIMESTAMP
        );`)

	return err
}

func createPgConnStr(cfg *config.PostgresCfg) string {

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s sslmode=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Sslmode,
		cfg.Dbname,
	)
}
