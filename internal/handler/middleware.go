package handler

import (
	"net/http"

	"github.com/cubny/cart/internal/auth"
	"github.com/cubny/cart/internal/ctxutil"
	"github.com/cubny/cart/internal/jsonerror"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Middleware defines information needed by authentication middleware.
type Middleware struct {
	auth AuthProvider
}

// NewMiddleware instantiates new authentication middleware.
func NewMiddleware(ap AuthProvider) *Middleware {
	return &Middleware{auth: ap}
}

// MiddlewareHandle is a method type that represents Middleware Handle function.
type MiddlewareHandle func(httprouter.Handle) httprouter.Handle

// MiddlewareChain represents helper struct that is able to wrap httprouter.Handle with chain.
type MiddlewareChain struct {
	middlewares []MiddlewareHandle
}

// With appends more handlers to the chain.
func (chain MiddlewareChain) With(handlers ...MiddlewareHandle) MiddlewareChain {
	chain.middlewares = append(chain.middlewares, handlers...)
	return chain
}

// Wrap the handler with MiddlewareChain.
func (chain MiddlewareChain) Wrap(handler httprouter.Handle) httprouter.Handle {
	result := handler

	for i := len(chain.middlewares) - 1; i >= 0; i-- {
		middleware := chain.middlewares[i]

		if result == nil {
			result = middleware(handler)
			continue
		}

		result = middleware(result)
	}

	return result
}

// Chain a bunch of chain.
func (middleware *Middleware) Chain(middlewares ...MiddlewareHandle) MiddlewareChain {
	return MiddlewareChain{middlewares: middlewares}
}

// Authorise checks whether client is auth to make the request.
func (middleware *Middleware) Authorise(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		key := auth.KeyFromClientRequest(r)

		if key == "" {
			_ = jsonerror.Unauthorised(w, "incorrect access_key")
			return
		}

		// Received an Key, verify it.
		accessKey, err := middleware.auth.VerifyKey(r.Context(), key)
		if err != nil {
			_ = jsonerror.Unauthorised(w, "incorrect access_key")
			return
		}

		ctx, err := ctxutil.SetUserAuthAccessKey(r.Context(), accessKey)
		if err != nil {
			log.Errorf("handler.middleware.auth: could not set user(%d) access key in context %s", accessKey.UserID, err)
			_ = jsonerror.InternalError(w, "")
			return
		}

		next(w, r.WithContext(ctx), ps)
	}
}

// ContentTypeJSON for the HTTP response.
func (middleware *Middleware) ContentTypeJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r, ps)
	}
}
