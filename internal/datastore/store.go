package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
)

type Store interface {
	// user
	CreateUser(id string, u *entities.User) error
	GetUserById(id string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	UpdateUser(id, name, username string) error
	DeleteUser(id string) error

	// product
	CreateProduct(id, sellerId string, p *entities.Product) error
	GetProductById(id string) (*entities.Product, error)
	UpdateProduct(id string, p *entities.Product) error
	UpdateStockProduct(id string, stock int) error
	DeleteProduct(id string) error
	GetProductSeller(id string) (string, error)
	ListProducts(q types.ListQueryValid, userId string) (*[]entities.Product, error)

	// bankAccount
	CreateBankAccount(id, sellerId string, b *entities.BankAccount) error
	GetBankAccount(id string) (*entities.BankAccount, error)
	ListBankAccount(id string) (*[]entities.BankAccount, error)
	DeleteBankAccount(id string) error
	UpdateBankAccount(id string, p *entities.BankAccount) error

	// transaction
	CreateTransaction(id, buyerId, sellerId, productId string, total float64, t *entities.Transaction) error
	GetTransaction(id string) (*TransactionReturn, error)
	ListTransaction(q types.ListQueryTransactionValid, userId string) (*[]TransactionReturn, error)
	UpdateStatusTransaction(id, status string) error
}

type TransactionReturn struct {
	Transaction entities.TransactionMinimal `json:"transaction"`
	Product     entities.ProductMinimal     `json:"product"`
	Seller      entities.UserMinimal        `json:"seller"`
	Buyer       entities.UserMinimal        `json:"buyer"`
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {

	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateTransaction(id, buyerId, sellerId, productId string, total float64, t *entities.Transaction) error {

	query := `INSERT INTO transactions(
    id,
    status,
    productId,
    buyerId,
    sellerId,
    quantity,
    notes,
    total
    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`

	_, err := s.db.Exec(
		query,
		id,
		"menunggu",
		productId,
		buyerId,
		sellerId,
		t.Quantity,
		t.Notes,
		total,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTransaction(id string) (*TransactionReturn, error) {
	var transaction TransactionReturn
	query := `
       SELECT 
            transactions.id,
            transactions.status,
            transactions.total,
            transactions.quantity,
            transactions.notes,
            transactions.createdAt,
            transactions.updatedAt,

            products.id,
            products.name,
            products.price,
            products.imageUrl,
            products.condition,
            products.tags,
            
            sellers.id,
            sellers.name,
            sellers.username,

            buyers.id,
            buyers.name,
            buyers.username
            
    FROM 
        transactions
    LEFT JOIN 
        products ON transactions.productId = products.id
    LEFT JOIN 
        users AS sellers ON transactions.sellerId = sellers.id
    LEFT JOIN 
        users AS buyers ON transactions.buyerId = buyers.id
    WHERE transactions.id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&transaction.Transaction.ID,
		&transaction.Transaction.Status,
		&transaction.Transaction.Total,
		&transaction.Transaction.Quantity,
		&transaction.Transaction.Notes,
		&transaction.Transaction.CreatedAt,
		&transaction.Transaction.UpdatedAt,

		&transaction.Product.ID,
		&transaction.Product.Name,
		&transaction.Product.Price,
		&transaction.Product.ImageUrl,
		&transaction.Product.Condition,
		&transaction.Product.Tags,

		&transaction.Seller.ID,
		&transaction.Seller.Name,
		&transaction.Seller.Username,

		&transaction.Buyer.ID,
		&transaction.Buyer.Name,
		&transaction.Buyer.Username,
	)

	if err != nil {
		return nil, err
	}

	switch {
	case err == sql.ErrNoRows:
		return &TransactionReturn{}, fmt.Errorf("Transaction did not exists")
	case err != nil:
		log.Println(err)
		return &TransactionReturn{}, fmt.Errorf("Something went wrong")
	default:
		return &transaction, nil
	}
}

func (s *Storage) ListTransaction(q types.ListQueryTransactionValid, userId string) (*[]TransactionReturn, error) {
	baseQuery, params := GenerateQueryListTransaction(q, userId)
	var returnTransaction []TransactionReturn
	rows, err := s.db.Query(baseQuery, params...)

	if err != nil {
		log.Println("err inside ListTransaction", err)
		return &[]TransactionReturn{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction TransactionReturn
		if err := rows.Scan(
			&transaction.Transaction.ID,
			&transaction.Transaction.Status,
			&transaction.Transaction.Total,
			&transaction.Transaction.Quantity,
			&transaction.Transaction.Notes,
			&transaction.Transaction.CreatedAt,
			&transaction.Transaction.UpdatedAt,

			&transaction.Product.ID,
			&transaction.Product.Name,
			&transaction.Product.Price,
			&transaction.Product.ImageUrl,
			&transaction.Product.Condition,
			&transaction.Product.Tags,
			&transaction.Product.Descriptions,

			&transaction.Seller.ID,
			&transaction.Seller.Name,
			&transaction.Seller.Username,

			&transaction.Buyer.ID,
			&transaction.Buyer.Name,
			&transaction.Buyer.Username,
		); err != nil {
			return &[]TransactionReturn{}, nil
		}

		returnTransaction = append(returnTransaction, transaction)
	}

	if len(returnTransaction) == 0 {
		return &[]TransactionReturn{}, nil
	}

	return &returnTransaction, nil
}

func (s *Storage) UpdateStatusTransaction(id, status string) error {
	query := `
    UPDATE transactions 
    SET status = $1,
        updatedAt = NOW()
    WHERE id = $2;
    `

	_, err := s.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateBankAccount(id string, b *entities.BankAccount) error {

	query := `UPDATE bankAccounts 
    SET bankName = $1,
        accountName = $2,
        accountNumber = $3,
        updatedAt = NOW() 
    WHERE id = $4;`

	_, err := s.db.Exec(query, b.BankName, b.AccountName, b.AccountNumber, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteBankAccount(id string) error {

	query := `DELETE FROM bankAccounts WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListBankAccount(id string) (*[]entities.BankAccount, error) {

	var returnBankAcc []entities.BankAccount
	query := `
    SELECT 
        id,
        bankName,
        accountName,
        accountNumber,
        sellerId,
        createdAt,
        updatedAt,
        deletedAt 
    FROM bankAccounts 
    WHERE sellerId = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Println("err inside ListBankAccount", err)
		return &[]entities.BankAccount{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var bankAcc entities.BankAccount
		if err := rows.Scan(
			&bankAcc.Id,
			&bankAcc.BankName,
			&bankAcc.AccountName,
			&bankAcc.AccountNumber,
			&bankAcc.SellerId,
			&bankAcc.CreatedAt,
			&bankAcc.UpdatedAt,
			&bankAcc.DeletedAt); err != nil {

			return &[]entities.BankAccount{}, nil
		}

		returnBankAcc = append(returnBankAcc, bankAcc)
	}

	return &returnBankAcc, nil
}

func (s *Storage) GetBankAccount(id string) (*entities.BankAccount, error) {

	var bankAccount entities.BankAccount
	query := `
    SELECT 
        id,
        bankName,
        accountName,
        accountNumber,
        sellerId,
        createdAt,
        updatedAt,
        deletedAt 
    FROM bankAccounts WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&bankAccount.Id,
		&bankAccount.BankName,
		&bankAccount.AccountName,
		&bankAccount.AccountNumber,
		&bankAccount.SellerId,
		&bankAccount.CreatedAt,
		&bankAccount.UpdatedAt,
		&bankAccount.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	switch {
	case err == sql.ErrNoRows:
		return &entities.BankAccount{}, fmt.Errorf("Bank Account did not exists")
	case err != nil:
		log.Println(err)
		return &entities.BankAccount{}, fmt.Errorf("Something went wrong")
	default:
		return &bankAccount, nil
	}

}

func (s *Storage) CreateBankAccount(id, sellerId string, b *entities.BankAccount) error {
	query := `
    INSERT INTO bankAccounts (
        id,
        bankName,
        accountName,
        accountNumber,
        sellerId
    ) VALUES ($1,$2,$3,$4,$5)`
	_, err := s.db.Exec(query, id, b.BankName, b.AccountName, b.AccountNumber, sellerId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateUser(id string, u *entities.User) error {

	hashedPassword := helper.GenerateHash(u.HashPassword)

	_, err := s.db.Exec(`
        INSERT INTO users (
            id,
            username,
            name,
            hashPassword
        )VALUES ($1,$2,$3,$4)`,
		id,
		u.Username,
		u.Name,
		hashedPassword,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserById(id string) (*entities.User, error) {

	var user entities.User
	err := s.db.QueryRow(`
        SELECT 
            id,
            name,
            username,
            hashPassword,
            createdAt,
            updatedAt,
            deletedAt 
        FROM users 
        WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.HashPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return &entities.User{}, fmt.Errorf("User did not exists")
	case err != nil:
		log.Println(err)
		return &entities.User{}, fmt.Errorf("Something went wrong")
	default:
		return &user, nil
	}

}

func (s *Storage) GetUserByUsername(username string) (*entities.User, error) {

	var user entities.User
	err := s.db.QueryRow(`
        SELECT 
            id, 
            name,
            username,
            hashPassword,
            createdAt,
            updatedAt,
            deletedAt 
        FROM users 
        WHERE username = $1`, username).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.HashPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return &entities.User{}, fmt.Errorf("User did not exists")
	case err != nil:
		log.Println(err)
		return &entities.User{}, fmt.Errorf("Something went wrong")
	default:
		return &user, nil
	}

}

func (s *Storage) UpdateUser(id, username, name string) error {

	_, err := s.db.Exec(`
        UPDATE users 
        SET name = $1,
            username = $2,
            updatedAt = NOW() 
        WHERE id = $3;
        `, name, username, id)

	if err != nil {

		return err
	}

	return nil
}

func (s *Storage) DeleteUser(id string) error {

	res, err := s.db.Exec(`
        DELETE FROM users
        WHERE id = $1;
        `, id)

	if err != nil {
		return err
	}
	rowAffect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffect < 1 {
		return fmt.Errorf("User didnot exists")
	}

	return nil
}

func (s *Storage) CreateProduct(id, sellerId string, p *entities.Product) error {
	tagArray := "{" + helper.ArrayToString(p.Tags) + "}"
	_, err := s.db.Exec(`
        INSERT INTO products (
            id,
            name,
            price,
            imageUrl, 
            condition,
            tags,
            isPurchaseable, 
            sellerId,
            stock
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
        `,
		id,
		p.Name,
		p.Price,
		p.ImageUrl,
		p.Condition,
		tagArray,
		p.IsPurchaseable,
		sellerId,
		p.Stock,
		p.Descriptions,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListProducts(q types.ListQueryValid, userId string) (*[]entities.Product, error) {

	baseQuery, params := GenerateQueryListProduct(q, userId)
	var returnProducts []entities.Product
	rows, err := s.db.Query(baseQuery, params...)

	if err != nil {
		log.Println("err inside ListProducts", err)
		return &[]entities.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
			&product.Condition,
			&product.Tags,
			&product.IsPurchaseable,
			&product.SellerId,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
			&product.Descriptions,
		); err != nil {

			log.Println(err)
			return &[]entities.Product{}, nil
		}

		returnProducts = append(returnProducts, product)
	}

	if len(returnProducts) == 0 {
		return &[]entities.Product{}, err
	}

	return &returnProducts, nil
}

func (s *Storage) GetProductById(id string) (*entities.Product, error) {

	var product entities.Product

	err := s.db.QueryRow(`
        SELECT 
            id, 
            name,
            price, 
            imageUrl, 
            condition,
            tags,  
            isPurchaseable,
            sellerId,
            stock,
            descriptions,
            createdAt,
            updatedAt,
            deletedAt
        FROM products 
        WHERE id = $1`, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.Condition,
		&product.Tags,
		&product.IsPurchaseable,
		&product.SellerId,
		&product.Stock,
		&product.Descriptions,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return &entities.Product{}, fmt.Errorf("Product did not exists")
	case err != nil:
		log.Println(err)
		return &entities.Product{}, fmt.Errorf("Something went wrong")
	default:
		return &product, nil
	}

}

func (s *Storage) UpdateProduct(id string, p *entities.Product) error {

	tagArray := "{" + helper.ArrayToString(p.Tags) + "}"

	_, err := s.db.Exec(`
        UPDATE products
        SET name = $1,
            price = $2,
            imageUrl = $3,
            stock = $4,
            condition = $5,
            tags = $6,
            isPurchaseable = $7,
            updatedAt = NOW()
        WHERE id = $9`,
		p.Name,
		p.Price,
		p.ImageUrl,
		p.Stock,
		p.Condition,
		tagArray,
		p.IsPurchaseable,
		id)

	if err != nil {

		return err
	}

	return nil
}

func (s *Storage) DeleteProduct(id string) error {

	_, err := s.db.Exec(`DELETE FROM products WHERE id = $1;`, id)
	if err != nil {

		return err
	}

	return nil
}

func (s *Storage) GetProductSeller(id string) (string, error) {

	var sellerId string

	if err := s.db.QueryRow(`SELECT sellerId FROM products WHERE id = $1`, id).Scan(&sellerId); err != nil {

		return "", err
	}

	return sellerId, nil
}

func (s *Storage) UpdateStockProduct(id string, stock int) error {

	query := `UPDATE products 
    SET stock = $1,
        updatedAt = NOW()
    WHERE id = $2`
	_, err := s.db.Exec(query, stock, id)
	if err != nil {
		return err
	}

	return nil
}

func GenerateQueryListTransaction(q types.ListQueryTransactionValid, userId string) (string, []interface{}) {

	baseQuery := `
    SELECT 
        transactions.id,
        transactions.status,
        transactions.total,
        transactions.quantity,
        transactions.notes,
        transactions.createdAt,
        transactions.updatedAt,

        products.id,
        products.name,
        products.price,
        products.imageUrl,
        products.condition,
        products.tags,
        products.descriptions,
        
        sellers.id,
        sellers.name,
        sellers.username,

        buyers.id,
        buyers.name,
        buyers.username
    FROM transactions 
    LEFT JOIN 
        products ON transactions.productId = products.id
    LEFT JOIN 
        users AS sellers ON transactions.sellerId = sellers.id
    LEFT JOIN 
        users AS buyers ON transactions.buyerId = buyers.id
    WHERE `

	queryIndex := 1
	var params []interface{}

	if q.Seller == true {
		baseQuery += `(transactions.sellerId = $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, userId)
		queryIndex += 1
	} else {
		baseQuery += `(transactions.buyerId = $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, userId)
		queryIndex += 1
	}

	if q.Search != "" {
		baseQuery += `(transactions.id = $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, q.Search)
		queryIndex += 1
	}

	baseQuery = baseQuery[:len(baseQuery)-5]
	baseQuery += ` ORDER BY ` + q.Order

	if q.Sort == "asc" {
		baseQuery += ` ASC `
	} else {
		baseQuery += ` DESC `
	}

	baseQuery += ` LIMIT $` + strconv.Itoa(queryIndex)
	queryIndex += 1
	params = append(params, q.Limit)

	baseQuery += ` OFFSET $` + strconv.Itoa(queryIndex)
	params = append(params, q.Offset)
	baseQuery += `;`

	return baseQuery, params
}

func GenerateQueryListProduct(q types.ListQueryValid, userId string) (string, []interface{}) {

	baseQuery := `
    SELECT 
        id,
        name,
        price,
        imageUrl,
        condition,
        tags,
        isPurchaseable,
        sellerId,
        stock,
        createdAt,
        updatedAt,
        deletedAt,
        descriptions 
    FROM products WHERE `
	queryIndex := 1
	var params []interface{}

	if q.UserOnly == "true" && userId != "" {
		baseQuery += `(sellerId = $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, userId)
		queryIndex += 1
	}

	baseQuery += `(condition = $` + strconv.Itoa(queryIndex) + `) AND `
	params = append(params, q.Condition)
	queryIndex += 1

	if len(q.Tags) > 0 {
		baseQuery += `(tags IN($` + strconv.Itoa(queryIndex) + `)) AND `
		params = append(params, "{"+helper.ArrayToString(q.Tags)+"}")
		queryIndex += 1
	}

	if q.ShowEmptyStock == "false" {
		baseQuery += `(stock > 0) AND `
	} else {
		baseQuery += `(stock = 0) AND `
	}

	baseQuery += `(price >= $` + strconv.Itoa(queryIndex) + `) AND `
	params = append(params, q.MinPrice)
	queryIndex += 1

	if q.MaxPrice > 0 {
		baseQuery += `(price <= $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, q.MaxPrice)
		queryIndex += 1
	}

	if q.Search != "" {
		baseQuery += `(name LIKE $` + strconv.Itoa(queryIndex) + `) AND `
		params = append(params, "%"+q.Search+"%")
		queryIndex += 1
	}

	baseQuery = baseQuery[:len(baseQuery)-5]

	baseQuery += ` ORDER BY ` + q.Order

	if q.Sort == "asc" {
		baseQuery += ` ASC `
	} else {
		baseQuery += ` DESC `
	}

	baseQuery += ` LIMIT $` + strconv.Itoa(queryIndex)
	queryIndex += 1
	params = append(params, q.Limit)

	baseQuery += ` OFFSET $` + strconv.Itoa(queryIndex) + `;`
	params = append(params, q.Offset)
	baseQuery += `;`

	return baseQuery, params
}
