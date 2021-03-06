package address

import (
	"net/http"
	"testing"
	"time"

	"github.com/Melon-Network-Inc/account-service/internal/auth"
	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/internal/test"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Address{
		{
			ID:        "123",
			Name:      "address123",
			Owner:     "testaccount",
			Pubkey:    "9ZNTfG4NyQgxy2SWjSiQoUyBPEvXT2xo7fKc5hPYYJ7b",
			Currency:  "SOL",
			IsPrimary: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{
			Name:         "get all",
			Method:       "GET",
			URL:          "/addresses",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `*"total_count":1*`,
		},
		{
			Name:         "get 123",
			Method:       "GET",
			URL:          "/addresses/123",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `*address123*`,
		},
		{
			Name:         "get unknown",
			Method:       "GET",
			URL:          "/addresses/1234",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusNotFound,
			WantResponse: "",
		},
		{
			Name:         "create ok",
			Method:       "POST",
			URL:          "/addresses",
			Body:         `{"name":"test"}`,
			Header:       header,
			WantStatus:   http.StatusCreated,
			WantResponse: "*test*",
		},
		{
			Name:         "create ok count",
			Method:       "GET",
			URL:          "/addresses",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `*"total_count":2*`,
		},
		{
			Name:         "create auth error",
			Method:       "POST",
			URL:          "/addresses",
			Body:         `{"name":"test"}`,
			Header:       nil,
			WantStatus:   http.StatusUnauthorized,
			WantResponse: "",
		},
		{
			Name:         "create input error",
			Method:       "POST",
			URL:          "/addresses",
			Body:         `"name":"test"}`,
			Header:       header,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "",
		},
		{
			Name:         "update ok",
			Method:       "PUT",
			URL:          "/addresses/123",
			Body:         `{"name":"addressxyz"}`,
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: "*addressxyz*",
		},
		{
			Name:         "update verify",
			Method:       "GET",
			URL:          "/addresses/123",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `*addressxyz*`,
		},
		{
			Name:         "update auth error",
			Method:       "PUT",
			URL:          "/addresses/123",
			Body:         `{"name":"addressxyz"}`,
			Header:       nil,
			WantStatus:   http.StatusUnauthorized,
			WantResponse: "",
		},
		{
			Name:         "update input error",
			Method:       "PUT",
			URL:          "/addresses/123",
			Body:         `"name":"addressxyz"}`,
			Header:       header,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "",
		},
		{
			Name:         "delete ok",
			Method:       "DELETE",
			URL:          "/addresses/123",
			Body:         ``,
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: "*addressxyz*",
		},
		{
			Name:         "delete verify",
			Method:       "DELETE",
			URL:          "/addresses/123",
			Body:         ``,
			Header:       header,
			WantStatus:   http.StatusNotFound,
			WantResponse: "",
		},
		{
			Name:         "delete auth error",
			Method:       "DELETE",
			URL:          "/addresses/123",
			Body:         ``,
			Header:       nil,
			WantStatus:   http.StatusUnauthorized,
			WantResponse: "",
		},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
