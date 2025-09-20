package order

import "time"

type Order struct {
	ID          int64       `json:"id" db:"id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	Status      string      `json:"status" db:"status"`
	TotalAmount int64       `json:"total_amount" db:"total_amount"`
	Items       []OrderItem `json:"items" db:"-"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ID        int64 `json:"id" db:"id"`
	OrderID   int64 `json:"order_id" db:"order_id"`
	ProductID int64 `json:"product_id" db:"product_id"`
	Quantity  int   `json:"quantity" db:"quantity"`
	Price     int64 `json:"price" db:"price"`
}
