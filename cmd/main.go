package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/liju-github/internal/config"
	"github.com/liju-github/internal/controller"
	"github.com/liju-github/internal/middleware"
	"github.com/liju-github/internal/repository"
	"github.com/liju-github/internal/service"
)

func main() {
	password := url.QueryEscape("Qwerty@123")

	mongoURI := "mongodb+srv://admin:" + password + "@olx.ibbnf.mongodb.net/?retryWrites=true&w=majority&appName=olx"
	dbName := "olxDB"
	db, err := config.NewMongoDB(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer db.Client.Disconnect(context.TODO())

	// Repositories
	userRepo := repository.UserRepository{Collection: db.Database.Collection("users")}
	productRepo := repository.ProductRepository{Collection: db.Database.Collection("products")}

	userService := &service.UserService{UserRepo: userRepo}
	productService := &service.ProductService{ProductRepo: productRepo}

	// Controllers
	userController := &controller.UserController{
		UserService:    userService,
		ProductService: productService,
	}
	productController := &controller.ProductController{ProductService: productService}

	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	router.POST("/signup", userController.Signup)
	router.POST("/login", userController.Login)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(userRepo))

	authRoutes.POST("/addproduct", productController.AddProduct)
	authRoutes.GET("/getproducts", productController.GetAllProducts)
	authRoutes.GET("/allusers", userController.GetAllUsers)

	gracefulShutdown(router)

	log.Println("Server started on port 8080")
}

func gracefulShutdown(router *gin.Engine) {
	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Server stopped unexpectedly: %v", err)
		}
	}()

	// Listen for interrupt signals (Ctrl+C, kill, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	// Shutdown the server with a timeout.
	timeout := 5 * time.Second
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Server gracefully stopped")
}
