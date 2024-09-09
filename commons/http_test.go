package commons

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleServiceError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		wantsErr     bool
	}{
		{"No Error", nil, http.StatusOK, false},
		{"SQL No Rows Error", sql.ErrNoRows, http.StatusNotFound, false},
		{"Generic Error", errors.New("generic error"), http.StatusInternalServerError, false},
		{"Unknown Error", errors.New("unknown error"), http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			c := e.NewContext(req, response)
			err := HandleServiceError(c, tt.err)
			if err != nil && tt.wantsErr {
				t.Errorf("Expected error '%v', got '%v'", nil, err)
			}
			if response.Code != tt.expectedCode {
				t.Errorf("Expected response code '%d', got '%d'", tt.expectedCode, response.Code)
			}

		})
	}
}
