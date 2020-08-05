package cart

import (
	"errors"
	"time"
)

var ErrInvalidUserID = errors.New("userID is not valid")

// Cart holds the basic data of a shopping cart
type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Price is a value type for price
type Price float64

// Item represents a fixed number of a single product in the shopping cart
type Item struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	CartID    int64 `json:"cart_id"`
	Quantity  int64 `json:"quantity"`

	// Price is the total price of the item, i.e. product's price * quantity
	Price     Price     `json:"price"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// NewCart creates a new cart
func NewCart(userID int64) (*Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	return &Cart{
		UserID: userID,
	}, nil
}
