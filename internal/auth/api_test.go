package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/Melon-Network-Inc/account-service/internal/errors"
	"github.com/Melon-Network-Inc/account-service/internal/test"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
)

type mockService struct{}

func (m mockService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "test" && password == "pass" {
		return "token-100", nil
	}
	return "", errors.Unauthorized("")
}

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	RegisterHandlers(router.Group(""), mockService{}, logger)

	tests := []test.APITestCase{
		{
			Name:         "success",
			Method:       "POST",
			URL:          "/login",
			Body:         `{"username":"test","password":"pass"}`,
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `{"token":"token-100"}`,
		},
		{
			Name:         "bad credential",
			Method:       "POST",
			URL:          "/login",
			Body:         `{"username":"test","password":"wrong pass"}`,
			Header:       nil,
			WantStatus:   http.StatusUnauthorized,
			WantResponse: "",
		},
		{
			Name:         "bad json",
			Method:       "POST",
			URL:          "/login",
			Body:         `"username":"test","password":"wrong pass"}`,
			Header:       nil,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "",
		},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
