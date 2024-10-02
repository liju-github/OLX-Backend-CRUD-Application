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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
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

type UserProfileResponse struct {
	Name     string          `json:"name"`
	ImageURL string          `json:"image_url"`
	Email    string          `json:"email"`
	Products []model.Product `json:"products"`
}

func (ctrl *UserController) GetSellerProfile(c *gin.Context) {
    email := c.Query("email")
    if email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
        return
    }

    sellerProfile, err := ctrl.UserService.GetUserByEmail(email)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
        return
    }

    products, err := ctrl.ProductService.GetAllProductsByUserEmail(email)
    if err != nil {
        log.Println("Failed to retrieve products for user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
        return
    }

    response := UserProfileResponse{
        Name:     sellerProfile.Name,
        ImageURL: sellerProfile.ImageURL,
        Email:    sellerProfile.Email,
        Products: products,
    }

    c.JSON(http.StatusOK, gin.H{
        "data": response,
    })
}
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// Retrieve the user email from the context
	userEmail, exists := c.Get("useremail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in the context"})
		return
	}

	// Fetch the user by email
	user, err := ctrl.UserService.GetUserByEmail(userEmail.(string))
	if err != nil {
		log.Println("Failed to retrieve user by email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Fetch the total number of products associated with the user
	products, err := ctrl.ProductService.GetAllProductsByUserEmail(user.Email)
	if err != nil {
		log.Println("Failed to retrieve products for user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// Create the response struct
	response := UserProfileResponse{
		Name:     user.Name,
		ImageURL: user.ImageURL,
		Email:    user.Email,
		Products: products,
	}

	c.JSON(http.StatusOK, gin.H{"profile": response})
}

type UpdateImageRequest struct {
	ImageUrl string `json:"image_url"` // Struct tag for JSON binding
}

func (controller *UserController) UpdateImage(c *gin.Context) {
    // Extract userEmail from context
    userEmail, exists := c.Get("useremail")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
        return
    }

    // Assert the type of userEmail to string
    email, ok := userEmail.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user email type"})
        return
    }
    
    // Initialize the request structure
    var req UpdateImageRequest

    // Bind the JSON payload to the request structure
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Check if image_url is provided
    if req.ImageUrl == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL is required"})
        return
    }

    // Call the service method, passing the Gin context as context.Context
    if err := controller.UserService.UpdateUserImage(c.Request.Context(), email, req.ImageUrl); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User image updated successfully."})
}
