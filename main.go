package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"project1/handler"
	"project1/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Connect to database using pgx driver
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if connection is established
	if err = db.Ping(); err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	fmt.Println("Connected to database")

	// Migrate database
	if _, err := Migrate(db); err != nil {
		fmt.Println("Error migrating database:", err)
		os.Exit(1)
	}

	// Create router using gin
	r := gin.Default()

	// Define routes

	// Public routes
	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout", handler.CheckoutOrder(db))

	r.POST("/api/v1/orders/:id/confirm", handler.ConfirmOrder(db))
	r.GET("/api/v1/orders/:id", handler.GetOrderById(db))

	// Admin routes
	r.POST("/admin/products", middleware.AdminOnly(), handler.CreateProduct(db))
	r.PUT("/admin/products/:id", middleware.AdminOnly(),  handler.UpdateProduct(db))
	r.DELETE("/admin/products/:id", middleware.AdminOnly(),  handler.DeleteProduct(db))
	

	// Start server
	server := http.Server{
		Addr: ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	defer server.Close()


	
}

