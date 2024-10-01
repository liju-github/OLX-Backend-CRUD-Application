package service

import (
    "github.com/liju-github/internal/model"
    "github.com/liju-github/internal/repository"
)

type ProductService struct {
    ProductRepo repository.ProductRepository
}

func (service *ProductService) AddProduct(product model.Product) error {
    return service.ProductRepo.AddProduct(product)
}

func (service *ProductService) GetAllProducts() ([]model.Product, error) {
    return service.ProductRepo.GetAllProducts()
}

func (service *ProductService)GetAllProductsByUserEmail(email string) ([]model.Product,error) {
    return service.ProductRepo.GetAllProductsByUserEmail(email)
}