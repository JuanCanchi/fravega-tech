package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fravega-tech/internal/domain"
	"fravega-tech/internal/handler"
	"fravega-tech/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProductRepository struct{}

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

func setupRouter() *gin.Engine {
	r := gin.Default()
	mockRepo := &MockProductRepository{}
	usecase := usecase.NewProductUsecase(mockRepo)
	handler.NewProductHandler(r, usecase)
	return r
}

func TestCreateProduct(t *testing.T) {
	router := setupRouter()

	product := domain.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Categories:  []string{"category1"},
		ImageURL:    "http://example.com/image.jpg",
	}

	body, _ := json.Marshal(product)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetProducts(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products?name=test", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUpdateProduct(t *testing.T) {
	router := setupRouter()

	product := domain.Product{
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       199.99,
		Categories:  []string{"category2"},
		ImageURL:    "http://example.com/new-image.jpg",
	}

	body, _ := json.Marshal(product)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/products/12345", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteProduct(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/products?ids=12345", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}
