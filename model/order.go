package model

import (
	"database/sql"
	"time"
)

type Checkout struct {
	Email string `json:"email"`
	Address string `json:"address"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	ID string `json:"id"`
	Quantity int32 `json:"quantity"`
}

type Order struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Address string `json:"address"`
	GrandTotal int64 `json:"grand_total"`
	Passcode *string `json:"passcode,omitempty"`
	PaidAt *time.Time `json:"paid_at,omitempty"`
	PaidBank *string `json:"paid_bank,omitempty"`
	PaidAccountNumber *string `json:"paid_account_number,omitempty"`
}

type OrderDetail struct {
	ID string `json:"id"`
	OrderID string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity int32 `json:"quantity"`
	Price int64 `json:"price"`
	Total int64 `json:"total"`
}


type OrderWithDetails struct {
	Order
	Details []OrderDetail `json:"details"`
}

type ConfirmOrder struct {
	Amount int64 `json:"amount" binding:"required"`
	Bank string `json:"bank" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	Passcode string `json:"passcode" binding:"required"`
}

func CreateOrder(db *sql.DB, order Order, details []OrderDetail) error {
	if db == nil {
		return errDBNil
	}
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO orders (id, email, address, passcode, grand_total) VALUES ($1, $2, $3, $4, $5);", order.ID, order.Email, order.Address, order.Passcode, order.GrandTotal)
	if err != nil {
		tx.Rollback()
		return err
	}

	queryDetails := "INSERT INTO order_details (id, order_id, product_id, quantity, price, total) VALUES ($1, $2, $3, $4, $5, $6);"
	
	for _, detail := range details {
		_, err = tx.Exec(queryDetails, detail.ID, detail.OrderID, detail.ProductID, detail.Quantity, detail.Price, detail.Total)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetOrderByID(db *sql.DB, id string) (Order, error) {
	if db == nil {
		return Order{}, errDBNil
	}
	var order Order
	query := "SELECT id, email, address, passcode, grand_total, paid_at, paid_bank, paid_account_number  FROM orders WHERE id = $1;"
	row := db.QueryRow(query, id)
	
	err := row.Scan(&order.ID, &order.Email, &order.Address, &order.Passcode, &order.GrandTotal, &order.PaidAt, &order.PaidBank, &order.PaidAccountNumber)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func UpdateOrderByID(db *sql.DB, id string, currentTime time.Time, bankName string, accountNumber string ) error {
	if db == nil {
		return errDBNil
	}
	query := "UPDATE orders SET paid_at = $1, paid_bank = $2, paid_account_number = $3 WHERE id = $4;"
	_, err := db.Exec(query, currentTime, bankName, accountNumber, id)
	if err != nil {
		return err
	}
	return nil
}


func GetOrderWithDetailsByID(db *sql.DB, id string) (OrderWithDetails, error) {
	if db == nil {
		return OrderWithDetails{}, errDBNil
	}

	order, err := GetOrderByID(db, id)
	if err != nil {
		return OrderWithDetails{}, err
	}

	query := "SELECT id, order_id, product_id, quantity, price, total FROM order_details WHERE order_id = $1;"
	rows, err := db.Query(query, id)
	if err != nil {
		return OrderWithDetails{}, err
	}

	var details []OrderDetail

	for rows.Next() {
		var detail OrderDetail
		err := rows.Scan(&detail.ID, &detail.OrderID, &detail.ProductID, &detail.Quantity, &detail.Price, &detail.Total)
		if err != nil {
			return OrderWithDetails{}, err
		}
		details = append(details, detail)
	}

	return OrderWithDetails{Order: order, Details: details}, nil


}