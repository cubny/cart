package service

import (
	"context"
	"errors"

	"github.com/cubny/cart"
	"github.com/cubny/cart/internal/storage"
)

var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrItemNotFound         = errors.New("item not found")
	ErrProductAlreadyInCart = errors.New("product is already in the cart")
)

// Service contains all the business logic of the shopping cart
type Service struct {
	storage Storage
}

// Storage provides the methods to CRUD resources in database
type Storage interface {
	CreateCart(ctx context.Context, cart *cart.Cart) error
	GetCart(ctx context.Context, userID, cartID int64) (*cart.Cart, error)
	FindItemByProductID(ctx context.Context, cartID, productID int64) (*cart.Item, error)
	CreateItem(ctx context.Context, item *cart.Item) error
	GetItem(ctx context.Context, itemID int64) (*cart.Item, error)
	RemoveItem(ctx context.Context, itemID int64) error
	RemoveItemsByCartID(ctx context.Context, cartID int64) error
	Close() error
}

// New creates a new Service
func New(db Storage) (*Service, error) {
	return &Service{storage: db}, nil
}

// CreateCart creates and persists a new cart for the given user
// this is only for demonstration, in real life, it should first check
// if the user has a open cart already for that we would need to mark the
// cart as closed when it's converted to order
func (s *Service) CreateCart(ctx context.Context, userID int64) (*cart.Cart, error) {
	cart, err := cart.NewCart(userID)
	if err != nil {
		return nil, err
	}

	if err := s.storage.CreateCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// AddItem, adds a product to the user's cart, it first checks if the cart belongs
// to the user. it then checks if the product is already added to the cart, if the
// product was already added it returns error, if not it adds the item to the cart
func (s *Service) AddItem(ctx context.Context, userID int64, item *cart.Item) error {
	// check the ownership of the cart
	_, err := s.storage.GetCart(ctx, userID, item.CartID)
	switch {
	case err == storage.ErrRecordNotFound:
		return ErrCartNotFound
	case err != nil:
		return err
	}

	// check if the product already exists in the cart
	t, err := s.storage.FindItemByProductID(ctx, item.CartID, item.ProductID)
	switch {
	case err == storage.ErrRecordNotFound:
	case err != nil:
		return err
	case t != nil:
		return ErrProductAlreadyInCart
	}

	// persist the item in the storage
	if err := s.storage.CreateItem(ctx, item); err != nil {
		return err
	}

	return nil
}

// RemoveItem, removes an item from the cart
// it first checks if the cart belongs to the user and then removes the item
func (s *Service) RemoveItem(ctx context.Context, userID, itemID int64) error {
	item, err := s.storage.GetItem(ctx, itemID)
	switch {
	case err == storage.ErrRecordNotFound:
		return ErrItemNotFound
	case err != nil:
		return err
	}

	// check the ownership of the cart
	_, err = s.storage.GetCart(ctx, userID, item.CartID)
	switch {
	case err == storage.ErrRecordNotFound:
		return ErrCartNotFound
	case err != nil:
		return err
	}

	return s.storage.RemoveItem(ctx, item.ID)
}

// EmptyCart remove all items of a cart
// it first checks the ownership of the cart and then delete all items
func (s *Service) EmptyCart(ctx context.Context, userID, cartID int64) error {
	// check the ownership of the cart
	_, err := s.storage.GetCart(ctx, userID, cartID)
	switch {
	case err == storage.ErrRecordNotFound:
		return ErrCartNotFound
	case err != nil:
		return err
	}

	return s.storage.RemoveItemsByCartID(ctx, cartID)
}

// CartDetails collects all the data about a cart
func (s *Service) CartDetails(ctx context.Context, cartID int) (*cart.Cart, []cart.Item, error) {
	// TODO implement this
	return nil, nil, nil
}
