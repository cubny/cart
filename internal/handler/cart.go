package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cubny/cart"
	"github.com/cubny/cart/internal/ctxutil"
	"github.com/cubny/cart/internal/jsonerror"
	"github.com/cubny/cart/internal/service"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// createCart is the handler for
// POST /carts/
func (h *Handler) createCart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accessKey, err := ctxutil.GetUserAuthAccessKey(r.Context())
	if err != nil {
		_ = jsonerror.InternalError(w, "cannot retrieve user from accessKey")
		return
	}

	c, err := h.service.CreateCart(r.Context(), accessKey.UserID)
	switch {
	case err == cart.ErrInvalidUserID:
		_ = jsonerror.InvalidParams(w, "user is invalid")
		return
	case err != nil:
		log.WithError(err).Errorf("createCart: %s", err)
		_ = jsonerror.InternalError(w, "cannot create cart")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		_ = jsonerror.InternalError(w, "cannot encode response")
		return
	}
}

// addItem is the handler for
// POST /cart/:cartID/items
func (h *Handler) addItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accessKey, err := ctxutil.GetUserAuthAccessKey(r.Context())
	if err != nil {
		_ = jsonerror.InternalError(w, "cannot retrieve user from accessKey")
		return
	}

	cartID, err := strconv.Atoi(p.ByName("cartID"))
	if err != nil {
		_ = jsonerror.InvalidParams(w, "cart_id param is not a valid number")
		return
	}

	itemReq := &struct {
		ProductID int64   `json:"product_id"`
		Price     float64 `json:"price"`
		Quantity  int64   `json:"quantity"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(itemReq); err != nil {
		_ = jsonerror.BadRequest(w, "body has invalid json format")
		return
	}

	item := &cart.Item{
		ProductID: itemReq.ProductID,
		CartID:    int64(cartID),
		Quantity:  itemReq.Quantity,
		Price:     cart.Price(itemReq.Price),
	}

	err = h.service.AddItem(r.Context(), accessKey.UserID, item)
	switch {
	case err == service.ErrCartNotFound:
		_ = jsonerror.NotFound(w, "cart does not exist")
		return
	case err == service.ErrProductAlreadyInCart:
		_ = jsonerror.BadRequest(w, "an item with the same product exists in the cart")
		return
	case err != nil:
		_ = jsonerror.InternalError(w, "could not add item to cart")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		_ = jsonerror.InternalError(w, "cannot decode item")
		return
	}
}

// removeItem is the handler for
// DELETE /items/:itemID
func (h *Handler) removeItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accessKey, err := ctxutil.GetUserAuthAccessKey(r.Context())
	if err != nil {
		_ = jsonerror.InternalError(w, "cannot retrieve user from accessKey")
		return
	}

	itemID, err := strconv.Atoi(p.ByName("itemID"))
	if err != nil {
		_ = jsonerror.InvalidParams(w, "item_id param is not a valid number")
		return
	}

	err = h.service.RemoveItem(r.Context(), accessKey.UserID, int64(itemID))
	switch {
	case err == service.ErrItemNotFound:
		_ = jsonerror.NotFound(w, "item does not exist")
		return
	case err == service.ErrCartNotFound:
		_ = jsonerror.NotFound(w, "cart does not exist")
		return
	case err != nil:
		log.WithError(err).Errorf("removeItem: service could not remove item")
		_ = jsonerror.InternalError(w, "could not remove item")
	}

	w.WriteHeader(http.StatusNoContent)
}

// emptyCart is the handler for
// DELETE /carts/:cartID/items
func (h *Handler) emptyCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// to empty a cart there are a couple of other options such as PUT with a
	// desired status of the cart e.g. items:{} or status:empty or POST a command
	// like /carts/:cartID/empty (which does not comply with RESTful verbs fully,
	// although it is fine) but I find the following more readable and continent
	// DELETE all items of this resource
	accessKey, err := ctxutil.GetUserAuthAccessKey(r.Context())
	if err != nil {
		_ = jsonerror.InternalError(w, "cannot retrieve user from accessKey")
		return
	}

	cartID, err := strconv.Atoi(p.ByName("cartID"))
	if err != nil {
		_ = jsonerror.InvalidParams(w, "cart_id param is not a valid number")
		return
	}

	err = h.service.EmptyCart(r.Context(), accessKey.UserID, int64(cartID))
	switch {
	case err == service.ErrCartNotFound:
		_ = jsonerror.NotFound(w, "cart does not exist")
		return
	case err != nil:
		log.WithError(err).Errorf("emptyCart: service could not empty cart")
		_ = jsonerror.InternalError(w, "could not empty cart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
