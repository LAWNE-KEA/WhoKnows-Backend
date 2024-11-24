package handlers

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "whoKnows/database"
    "whoKnows/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// Mock services and database connection
func init() {
    var err error
    database.Connection, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    database.Connection.AutoMigrate(&models.User{}, &models.Token{})
}

func TestLogoutHandler(t *testing.T) {
    tests := []struct {
        name           string
        authHeader     string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "Missing Authorization Header",
            authHeader:     "",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   "Authorization header missing",
        },
        {
            name:           "Invalid Token Prefix",
            authHeader:     "InvalidPrefix token",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   "Invalid token",
        },
        {
            name:           "Valid Token",
            authHeader:     "Bearer validtoken",
            expectedStatus: http.StatusOK,
            expectedBody:   `"status":"success"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest("POST", "/logout", nil)
            if err != nil {
                t.Fatal(err)
            }
            req.Header.Set("Authorization", tt.authHeader)

            rr := httptest.NewRecorder()
            handler := http.HandlerFunc(LogoutHandler)
            handler.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }

            if !strings.Contains(rr.Body.String(), tt.expectedBody) {
                t.Errorf("handler returned unexpected body: got %v want %v",
                    rr.Body.String(), tt.expectedBody)
            }
        })
    }
}