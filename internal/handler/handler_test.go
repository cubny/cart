//go:generate mockgen -package handler -source handler.go -destination handler_mocks_test.go
package handler_test

import (
	"testing"

	"github.com/cubny/cart/internal/handler"
	"github.com/cubny/cart/internal/tests"

	"github.com/stretchr/testify/assert"
)

func execHTTPTestCases(t *testing.T, sp handler.ServiceProvider, ap handler.AuthProvider, tcs []tests.TestCase) {
	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			handler, err := handler.New(sp, ap)
			assert.Nil(t, err)
			tests.HandlerTest(t, handler, &tc)
		})
	}
}
