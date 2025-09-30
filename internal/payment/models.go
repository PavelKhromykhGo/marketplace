package payment

type Intent struct {
	ID           int64  `json:"id" db:"id"`
	OrderID      int64  `json:"order_id" db:"order_id"`
	Amount       int64  `json:"amount" db:"amount"`
	Status       string `json:"status" db:"status"`
	ClientSecret string `json:"client_secret" db:"client_secret"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}
