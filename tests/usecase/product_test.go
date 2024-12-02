package usecase

import (
	"errors"
	"fravega-tech/internal/domain"
	usecase2 "fravega-tech/internal/usecase"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockProductRepository struct {
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	if product.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (m *MockProductRepository) GetAll(name, category string) ([]domain.Product, error) {
	return []domain.Product{{Name: "Product 1"}}, nil
}

func (m *MockProductRepository) Update(id string, product *domain.Product) error {
	if id == "" {
		return errors.New("invalid product id")
	}
	return nil
}

func (m *MockProductRepository) Delete(id string) error {
	if id == "" {
		return errors.New("invalid product id")
	}
	return nil
}

func (m *MockProductRepository) DeleteMany(ids []string) error {
	if len(ids) == 0 {
		return errors.New("no product ids provided")
	}
	return nil
}

func TestCreateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	usecase := usecase2.NewProductUsecase(mockRepo)

	product := &domain.Product{Name: "Test Product"}
	err := usecase.CreateProduct(product)

	assert.Nil(t, err)
}

func TestCreateProduct_InvalidName(t *testing.T) {
	mockRepo := &MockProductRepository{}
	usecase := usecase2.NewProductUsecase(mockRepo)

	product := &domain.Product{}
	err := usecase.CreateProduct(product)

	assert.NotNil(t, err)
	assert.Equal(t, "name is required", err.Error())
}

func TestGetProducts(t *testing.T) {
	mockRepo := &MockProductRepository{}
	usecase := usecase2.NewProductUsecase(mockRepo)

	products, err := usecase.GetProducts("Product 1", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(products))
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	usecase := usecase2.NewProductUsecase(mockRepo)

	product := &domain.Product{Name: "Updated Product"}
	err := usecase.UpdateProduct("12345", product)

	assert.Nil(t, err)
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := &MockProductRepository{}
	usecase := usecase2.NewProductUsecase(mockRepo)

	err := usecase.DeleteProducts([]string{"12345"})
	assert.Nil(t, err)
}
