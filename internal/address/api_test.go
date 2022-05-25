package address

import (
	"net/http"
	"testing"
	"time"
	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/internal/auth"
	"github.com/Melon-Network-Inc/account-service/internal/test"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Address{
		{"123", "address123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/addresses", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/addresses/123", "", nil, http.StatusOK, `*address123*`},
		{"get unknown", "GET", "/addresses/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/addresses", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/addresses", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/addresses", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/addresses", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/addresses/123", `{"name":"addressxyz"}`, header, http.StatusOK, "*addressxyz*"},
		{"update verify", "GET", "/addresses/123", "", nil, http.StatusOK, `*addressxyz*`},
		{"update auth error", "PUT", "/addresses/123", `{"name":"addressxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/addresses/123", `"name":"addressxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/addresses/123", ``, header, http.StatusOK, "*addressxyz*"},
		{"delete verify", "DELETE", "/addresses/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/addresses/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
