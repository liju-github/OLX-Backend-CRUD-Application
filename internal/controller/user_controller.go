package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/liju-github/internal/auth"
	"github.com/liju-github/internal/model"
	"github.com/liju-github/internal/service"
)

var validate = validator.New()

type UserController struct {
    UserService    *service.UserService
    ProductService *service.ProductService
}

func (ctrl *UserController) Signup(c *gin.Context) {
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        log.Println("Error binding JSON in Signup: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, unable to parse request body"})
        return
    }

            fmt.Println(user)

    if err := validate.Struct(user); err != nil {
        log.Println("Validation error in Signup: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
        return
    }

    if err := ctrl.UserService.RegisterUser(user); err != nil {
        log.Println("User registration failed: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user","details":err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (ctrl *UserController) Login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required"`
    }

    if err := c.ShouldBindJSON(&credentials); err != nil {
        log.Println("Error binding JSON in Login: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, unable to parse request body"})
        return
    }

    if err := validate.Struct(credentials); err != nil {
        log.Println("Validation error in Login: ", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
        return
    }

    user, err := ctrl.UserService.Login(credentials.Email, credentials.Password)
    if err != nil {
        log.Println("Login failed: ", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := auth.GenerateJWT(user.Email)
    if err != nil {
        log.Println("JWT generation failed: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    products, err := ctrl.ProductService.GetAllProductsByUserEmail(user.Email)
    if err != nil {
        log.Println("Failed to fetch products for user: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products for the user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user":     user,
        "products": products,
        "token":    token,
    })
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
    users, err := ctrl.UserService.AllUsers()
    if err != nil {
        log.Println("Error fetching all users: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}
