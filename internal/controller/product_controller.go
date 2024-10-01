package controller

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/liju-github/internal/model"
    "github.com/liju-github/internal/service"
)

type ProductController struct {
    ProductService *service.ProductService
}

func (ctrl *ProductController) AddProduct(c *gin.Context) {
    var product model.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        log.Println("Error binding JSON in AddProduct: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    // Validate the product fields
    if err := validate.Struct(product); err != nil {
        log.Println("Validation error in AddProduct: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

    userEmail, exists := c.Get("useremail")
    if !exists {
        log.Println("User email missing from context in AddProduct")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    email, ok := userEmail.(string)
    if !ok {
        log.Println("Invalid user email type in context")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    product.Email = email

    if err := ctrl.ProductService.AddProduct(product); err != nil {
        log.Println("Failed to add product in AddProduct: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func (ctrl *ProductController) GetProduct(c *gin.Context) {
    id := c.Param("id")

    product, err := ctrl.ProductService.ProductRepo.GetProductByID(id)
    if err != nil {
        log.Println("Failed to find product in GetProduct: ", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"product": product})
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
    products, err := ctrl.ProductService.ProductRepo.GetAllProducts()
    if err != nil {
        log.Println("Failed to fetch products in GetAllProducts: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"products": products})
}
