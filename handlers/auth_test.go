package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/auth"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"inventory-service-go/item"
	"inventory-service-go/person"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorize(t *testing.T) {
	tests := []struct {
		name               string
		credentials        *auth.Credentials
		expectedHTTPStatus int
		expectedResult     TokenCredentials
	}{
		{
			name:               "validCredentials",
			credentials:        &auth.Credentials{ClientId: "foo", ClientSecret: "bar"},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name:               "invalidCredentials",
			credentials:        &auth.Credentials{ClientId: "invalid", ClientSecret: "invalid"},
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "emptyCredentials",
			credentials:        &auth.Credentials{},
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "invalid request body",
			credentials:        nil,
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}
	controller := gomock.NewController(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock echo.Context
			e := echo.New()
			var req *http.Request
			if tc.expectedHTTPStatus == http.StatusBadRequest && tc.credentials == nil {
				invalidPayload := []byte(`not json`)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidPayload))
			} else {
				b, _ := json.Marshal(tc.credentials)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			}
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Mock ApplicationContext
			mockPersonService := person.NewMockPersonService(controller)
			mockItemService := item.NewMockItemService(controller)
			mockInvoiceService := invoice.NewMockInvoiceService(controller)
			mockContext := context.MockApplicationContext(mockPersonService, mockItemService, mockInvoiceService)

			// Test function
			if assert.NoError(t, Authorize(mockContext)(c)) {
				assert.Equal(t, tc.expectedHTTPStatus, rec.Code)
				if rec.Code == http.StatusOK {
					var tokenCredentials TokenCredentials
					_ = json.Unmarshal(rec.Body.Bytes(), &tokenCredentials)
					token := tokenCredentials.Token
					assert.NotEmpty(t, token)
					// validate JWT token
					decodedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
						return []byte("dummy_secret"), nil
					})
					assert.NotNil(t, decodedToken)
					assert.True(t, decodedToken.Valid)
				}
			}
		})
	}
}
