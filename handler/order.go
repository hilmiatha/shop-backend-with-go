package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"project1/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			log.Printf("Error binding product: %v", err)
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}


		ids := []string{}
		qtyMap := make(map[string]int32)
		for _, product := range checkoutOrder.Products {
			ids = append(ids, product.ID)
			qtyMap[product.ID] = product.Quantity
		}
		
		products, err := model.SelectProducts(db, ids)
		if err != nil {
			log.Printf("Error selecting products: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		//make password
		passcode := generatePasscode(5)

		// hash password
		hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		hashedPasscodeStr := string(hashedPasscode)

		order := model.Order{
			ID: uuid.NewString(),
			Email: checkoutOrder.Email,
			Address: checkoutOrder.Address,
			Passcode: &hashedPasscodeStr,
			GrandTotal: 0,
		}

		details:= []model.OrderDetail{}

		for _, product := range products {
			qty := qtyMap[product.ID]
			total := product.Price * int64(qty)
			order.GrandTotal += total
			details = append(details, model.OrderDetail{
				ID: uuid.NewString(),
				OrderID: order.ID,
				ProductID: product.ID,
				Quantity: qty,
				Price: product.Price,
				Total: total,
			})
		}

		err = model.CreateOrder(db, order, details)
		if err != nil {
			log.Printf("Error creating order: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		orderWithDetails := model.OrderWithDetails{
			Order: order,
			Details: details,
		}
		orderWithDetails.Order.Passcode = &passcode

		c.JSON(200, orderWithDetails)


	}
}

func generatePasscode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	passcode := make([]byte, length)
	for i := range passcode {
		passcode[i] = charset[randomGenerator.Intn(len(charset))]
	}
	return string(passcode)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		// bind id and request body to confirm order
		id := c.Param("id")

		var confirmOrder model.ConfirmOrder
		if err:= c.BindJSON(&confirmOrder); err != nil {
			log.Printf("Error binding product: %v", err)
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		// get order
		order, err := model.GetOrderByID(db, id)
		if err != nil {
			log.Printf("Error getting order: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// check passcode if exist
		if order.Passcode == nil {
			log.Println("Passcode not found")
			c.JSON(500, gin.H{"error": "passcode not found"})
			return
		}

		// compare passcode with confirm order passcode
		if bcryptErr := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirmOrder.Passcode)); bcryptErr != nil {
			log.Printf("Error comparing password: %v", bcryptErr)
			c.JSON(401, gin.H{"error": "passcode not match"})
			return
		}

		// check if order already paid
		if order.PaidAt != nil {
			log.Println("Order already paid")
			c.JSON(400, gin.H{"error": "order already paid"})
			return
		}

		// check if amount match with grand total order
		if order.GrandTotal != confirmOrder.Amount {
			log.Printf("Amount not match %d != %d", order.GrandTotal, confirmOrder.Amount)
			c.JSON(400, gin.H{"error": "amount not match"})
			return
		}

		current := time.Now()

		// update order
		if err := model.UpdateOrderByID(db, id, current, confirmOrder.Bank, confirmOrder.AccountNumber); err != nil {
			log.Printf("Error updating order: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}
		order.PaidAt = &current
		order.PaidBank = &confirmOrder.Bank
		order.PaidAccountNumber = &confirmOrder.AccountNumber
		// remove passcode because it's no longer needed after order is paid
		order.Passcode = nil

		c.JSON(200, order)
	}
}



func GetOrderById(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		id := c.Param("id")

		orderWithDetails, err := model.GetOrderWithDetailsByID(db, id)
		if err != nil {
			log.Printf("Error getting order: %v", err)
			c.JSON(500, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		c.JSON(200, orderWithDetails)
	}
}