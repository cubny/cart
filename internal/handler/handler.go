package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/cubny/cart"
	"github.com/cubny/cart/internal/auth"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	api500Count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "error_500_counter",
			Help:      "Counter of 500 responses of cart api",
		}, []string{"method", "reason"})
)

func init() {
	prometheus.MustRegister(api500Count)
}

// ServiceProvided contains all the business logic
type ServiceProvider interface {
	CreateCart(ctx context.Context, userID int64) (*cart.Cart, error)
	AddItem(ctx context.Context, userID int64, item *cart.Item) error
	RemoveItem(ctx context.Context, userID, itemID int64) error
	EmptyCart(ctx context.Context, userID, cartID int64) error
	//CartDetails(ctx context.Context, cartID int) (*cart.Cart, []cart.Item, error)
}

// AuthProvider provides the client to interact with the auth service
type AuthProvider interface {
	VerifyKey(ctx context.Context, accessKey string) (*auth.AccessKey, error)
}

// Handler handles http requests
type Handler struct {
	service ServiceProvider
	http.Handler
}

// New creates a new handler to handle http requests
func New(service ServiceProvider, authClient AuthProvider) (*Handler, error) {

	switch {
	case authClient == nil:
		return nil, errors.New("auth client is required")
	case service == nil:
		return nil, errors.New("service is required")
	}

	h := &Handler{
		service: service,
	}
	router := httprouter.New()

	middleware := NewMiddleware(authClient)
	chain := middleware.Chain(middleware.ContentTypeJSON, middleware.Authorise)

	router.GET("/health", h.health)
	router.POST("/carts", chain.Wrap(h.createCart))
	router.POST("/carts/:cartID/items", chain.Wrap(h.addItem))
	router.DELETE("/items/:itemID", chain.Wrap(h.removeItem))
	// to find out why I chose DELETE for emptying the cart read the comments in the handler
	router.DELETE("/carts/:cartID/items", chain.Wrap(h.emptyCart))

	h.Handler = router
	return h, nil
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
