package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"

	"github.com/GetterSethya/golangApiMarketplace/config"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/joho/godotenv"
)

type bankAccount struct {
	Id            string
	BankName      string
	AccountName   string
	AccountNumber int64
}

type userData struct {
	Id       string
	Name     string
	Username string
	Password string
}

type product struct {
	Id             string
	Name           string
	Price          float64
	ImageUrl       string
	Condition      string
	Tags           []string
	IsPurchaseable bool
	Stock          int
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		cfg.Postgres.User,
		cfg.Postgres.Dbname,
		cfg.Postgres.Sslmode,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	//user seed
	users := []userData{
		{
			Id:       "3e595902-9b50-49eb-96c9-178b1545bd80",
			Name:     "john",
			Username: "johndoe123",
			Password: "12345678",
		},
		{
			Id:       "93fcc1cc-68f4-4038-b3b9-3ec81ad0b4b4",
			Name:     "jane",
			Username: "janedoe123",
			Password: "12345678",
		},
	}

	for _, user := range users {
		if err := seedUser(db, user); err != nil {
			log.Fatal(err)
		}
	}

	//product seed
	products := []product{
		{
			Id:             "3f678471-b4b8-4757-ac36-1c44e458ad04",
			Name:           "Kue Nastar",
			Price:          15000,
			ImageUrl:       "example.com",
			Condition:      "new",
			Tags:           []string{"lebaran", "kue"},
			IsPurchaseable: true,
			Stock:          100,
		},
		{
			Id:             "fccfeaf7-0122-4920-a2db-41eda0487aa3",
			Name:           "Kue Putri Salju",
			Price:          20000,
			ImageUrl:       "example.com",
			Condition:      "new",
			Tags:           []string{"lebaran", "kue"},
			IsPurchaseable: true,
			Stock:          100,
		},
		{
			Id:             "71c98f33-8c45-4492-a81e-e588668da526",
			Name:           "Sendal swallow",
			Price:          12000,
			ImageUrl:       "example.com",
			Condition:      "new",
			Tags:           []string{"outfit", "keren"},
			IsPurchaseable: true,
			Stock:          100,
		},
		{
			Id:             "8c90dddf-176f-4ad1-ae79-3909531b70d9",
			Name:           "Ambatron",
			Price:          100000,
			ImageUrl:       "example.com",
			Condition:      "new",
			Tags:           []string{"pemimpin", "robot"},
			IsPurchaseable: true,
			Stock:          100,
		},
	}

	for _, product := range products {
		if err := seedProduct(db, product, users); err != nil {
			log.Fatal(err)
		}
	}

	//seed bank account
	bankAccs := []bankAccount{
		{
			Id:            "eb1b5f4a-dc1f-48aa-a994-4bef3ef839ad",
			BankName:      "BANK INI",
			AccountName:   "utama",
			AccountNumber: int64(rand.Intn(100000000)),
		},
		{
			Id:            "dd6e847e-2e26-4655-9cb6-9816174b5c9a",
			BankName:      "BANK ITU",
			AccountName:   "kedua",
			AccountNumber: int64(rand.Intn(100000000)),
		},
	}

	for _, bankAcc := range bankAccs {
		if err := seedBankAccount(db, bankAcc, users); err != nil {
			log.Fatal(err)
		}
	}

}

func seedBankAccount(db *sql.DB, b bankAccount, u []userData) error {

	_, err := db.Exec(`
        INSERT INTO bankAccounts (
        id,
        bankName,
        accountName,
        accountNumber,
        sellerId
        ) VALUES ($1,$2,$3,$4,$5)
        `,
		b.Id,
		b.BankName,
		b.AccountName,
		b.AccountNumber,
		u[rand.Intn(len(u))].Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func seedProduct(db *sql.DB, p product, u []userData) error {

	tagArray := "{" + helper.ArrayToString(p.Tags) + "}"
	_, err := db.Exec(`
        INSERT INTO products (
        id,
        name,
        price,
        imageUrl,
        condition,
        tags,
        isPurchaseable,
        sellerId,
        stock)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        `,
		p.Id,
		p.Name,
		p.Price,
		p.ImageUrl,
		p.Condition,
		tagArray,
		p.IsPurchaseable,
		u[rand.Intn(len(u))].Id,
		p.Stock,
	)

	if err != nil {
		return err
	}

	return nil
}

func seedUser(db *sql.DB, user userData) error {

	hash := helper.GenerateHash(user.Password)

	_, err := db.Exec(`
        INSERT INTO users (
        id,
        name,
        username,
        hashPassword)
        VALUES ($1,$2,$3,$4)
        `,
		user.Id,
		user.Name,
		user.Username,
		hash,
	)
	if err != nil {
		return err
	}

	return nil
}
