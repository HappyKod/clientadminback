package clientadmin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HappyKod/clientadminback"
	"github.com/HappyKod/clientadminback/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestWithLoginAndPassword1(t *testing.T) {
	loginAndPassword := clientadminback.LoginPassword{Email: "yudinsv@agatha-hub.ru", Password: "SAVA1973398sava"}
	clientAdmin := ClientAdmin{
		ServiceURL: "https://api-test.admin.agatha.pw/v1/auth/login",
	}
	err := WithLoginAndPassword(loginAndPassword)(&clientAdmin)
	if err != nil {
		t.Error(err)
	}
}

func TestWithLoginAndPassword(t *testing.T) {
	testCases := []struct {
		name           string
		loginPassword  clientadminback.LoginPassword
		clientInfo     models.ClientInfo
		expectedResult error
	}{
		{
			name: "successful login",
			loginPassword: clientadminback.LoginPassword{
				Email:    "test@example.com",
				Password: "password123",
			},
			clientInfo: models.ClientInfo{
				IsActive: true,
				Token:    "test-token",
			},
			expectedResult: nil,
		},
		{
			name: "failed login",
			loginPassword: clientadminback.LoginPassword{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			clientInfo: models.ClientInfo{
				IsActive: false,
			},
			expectedResult: ErrorAuthFailed,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				jsonData, _ := json.Marshal(tt.clientInfo)
				_, err := w.Write(jsonData)
				if err != nil {
					t.Error(err)
				}
			}))
			defer ts.Close()

			clientAdmin := &ClientAdmin{ServiceURL: ts.URL}

			option := WithLoginAndPassword(tt.loginPassword)
			result := option(clientAdmin)

			if tt.expectedResult == nil {
				assert.NoError(t, result)
			} else {
				assert.True(t, errors.Is(result, tt.expectedResult), "expected error: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestWithJWTToken(t *testing.T) {
	testCases := []struct {
		name           string
		JWT            string
		clientInfo     models.ClientInfo
		expectedResult error
	}{
		{
			name: "successful login",
			JWT:  "",
			clientInfo: models.ClientInfo{
				IsActive: true,
				Token:    "test-token",
			},
			expectedResult: nil,
		},
		{
			name: "failed login",
			JWT:  "",
			clientInfo: models.ClientInfo{
				IsActive: false,
			},
			expectedResult: ErrorAuthFailed,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				jsonData, _ := json.Marshal(tt.clientInfo)
				_, err := w.Write(jsonData)
				if err != nil {
					t.Error(err)
				}
			}))
			defer ts.Close()

			clientAdmin := &ClientAdmin{ServiceURL: ts.URL}
			option := WithJWTToken(tt.JWT)
			result := option(clientAdmin)

			if tt.expectedResult == nil {
				assert.NoError(t, result)
			} else {
				assert.True(t, errors.Is(result, tt.expectedResult), "expected error: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}