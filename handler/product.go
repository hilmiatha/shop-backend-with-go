package handler

import (
	"database/sql"
	"errors"
	"log"
	"project1/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: ambil dari db
		products, err := model.GetAllProducts(db)
		if err != nil {
			log.Printf("Error getting products: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}
		//TODO: beri response
		c.JSON(200, products)
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Baca id dari URL
		id := c.Param("id")
		
		//TODO: ambil dari db berdasarkan id
		product, err := model.GetProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
				log.Printf("Product not found: %v", err)
				c.JSON(404, gin.H{"error": "product not found"})
				return
			}
			log.Printf("Error getting product: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}
		//TODO: beri response 
		c.JSON(200, product)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		var product model.Product
		if err:= c.BindJSON(&product); err != nil {
			log.Printf("Error binding product: %v", err)
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		product.ID = uuid.NewString()

		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Error creating product: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}


func UpdateProduct(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		id := c.Param("id")

		var product model.Product
		if err:= c.BindJSON(&product); err != nil {
			log.Printf("Error binding product: %v", err)
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		productExist, err := model.GetProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
				log.Printf("Product not found: %v", err)
				c.JSON(404, gin.H{"error": "product not found"})
				return
			}
			log.Printf("Error getting product: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		if product.Name != "" {
			productExist.Name = product.Name
		}
		if product.Price != 0 {
			productExist.Price = product.Price
		}


		if err := model.UpdateProduct(db, productExist); err != nil {
			log.Printf("Error updating product: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		c.JSON(200, product)
	}
}


func DeleteProduct(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		id := c.Param("id")

		if err := model.DeleteProduct(db, id); err != nil {
			if errors.Is(err, sql.ErrNoRows){
				log.Printf("Product not found: %v", err)
				c.JSON(404, gin.H{"error": "product not found"})
				return
			}
			log.Printf("Error deleting product: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		c.JSON(204, nil)
	}
}