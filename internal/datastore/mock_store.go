package datastore

import (
	"github.com/GetterSethya/golangApiMarketplace/internal/entities"
	"github.com/GetterSethya/golangApiMarketplace/internal/types"
)

type MockStore struct {
}

func (m *MockStore) CreateBankAccount(id, sellerId string, b *entities.BankAccount) error {

	return nil
}

func (m *MockStore) GetBankAccount(id string) (*entities.BankAccount, error) {

	return &entities.BankAccount{}, nil
}

func (m *MockStore) ListBankAccount(id string) (*[]entities.BankAccount, error) {

	return &[]entities.BankAccount{}, nil
}

func (m *MockStore) DeleteBankAccount(id string) error {

	return nil
}

func (m *MockStore) UpdateBankAccount(id string, p *entities.BankAccount) error {

	return nil
}

func (m *MockStore) UpdateStockProduct(id string, stock int) error {

	return nil
}

func (m *MockStore) DeleteProduct(id string) error {

	return nil
}

func (m *MockStore) GetProductSeller(id string) (string, error) {

	return "75ea96d2-8077-48aa-aad6-a02fbd282f3c", nil
}

func (m *MockStore) ListProducts(q types.ListQueryValid, userId string) (*[]entities.Product, error) {

	return &[]entities.Product{}, nil
}

func (m *MockStore) CreateUser(id string, u *entities.User) error {

	return nil
}

func (m *MockStore) DeleteUser(id string) error {

	return nil
}

func (m *MockStore) UpdateUser(id, name, username string) error {

	return nil
}

func (m *MockStore) GetUserById(id string) (*entities.User, error) {

	return &entities.User{}, nil
}

func (m *MockStore) GetUserByUsername(username string) (*entities.User, error) {

	return &entities.User{}, nil
}

func (m *MockStore) CreateProduct(id, sellerId string, p *entities.Product) error {

	return nil
}

func (m *MockStore) GetProductById(id string) (*entities.Product, error) {

	return &entities.Product{}, nil
}

func (m *MockStore) SearchProduct(q string) (*entities.Product, error) {

	return &entities.Product{}, nil
}

func (m *MockStore) UpdateProduct(id string, p *entities.Product) error {

	return nil
}
