package usecase

import (
	"errors"
	"fravega-tech/internal/domain"
	"time"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	GetAll(name, category string) ([]domain.Product, error)
	Update(id string, product *domain.Product) error
	Delete(id string) error
	DeleteMany(ids []string) error
}

type ProductUsecase struct {
	repo ProductRepository
}

func NewProductUsecase(repo ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(product *domain.Product) error {
	if product.Price < 0 {
		return errors.New("price cannot be negative")
	}
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	return u.repo.Create(product)
}

func (u *ProductUsecase) GetProducts(name, category string) ([]domain.Product, error) {
	products, err := u.repo.GetAll(name, category)
	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []domain.Product{}
	}

	return products, nil
}

func (u *ProductUsecase) UpdateProduct(id string, product *domain.Product) error {
	now := time.Now()
	product.UpdatedAt = now

	product.ID = id

	return u.repo.Update(id, product)
}

func (u *ProductUsecase) DeleteProducts(ids []string) error {
	return u.repo.DeleteMany(ids)
}
