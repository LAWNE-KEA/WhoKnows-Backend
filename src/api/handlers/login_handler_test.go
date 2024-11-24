package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "whoKnows/database"
    "whoKnows/models"
    "whoKnows/security"
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
    database.Connection.AutoMigrate(&models.User{})
}

func TestLoginHandler(t *testing.T) {
    // Create a mock user
    hashedPassword, err := security.HashPassword("testpassword")
    if err != nil {
        t.Fatal(err)
    }
    user := models.User{
        Username: "testuser",
        Password: hashedPassword,
    }
    database.Connection.Create(&user)

    tests := []struct {
        name           string
        requestBody    LoginRequest
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "Valid login",
            requestBody: LoginRequest{
                Username: "testuser",
                Password: "testpassword",
            },
            expectedStatus: http.StatusOK,
            expectedBody:   `"status":"success"`,
        },
        {
            name: "Invalid login - missing username",
            requestBody: LoginRequest{
                Password: "testpassword",
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `"username and password are required"`,
        },
        {
            name: "Invalid login - wrong password",
            requestBody: LoginRequest{
                Username: "testuser",
                Password: "wrongpassword",
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   `"invalid username or password"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
            if err != nil {
                t.Fatal(err)
            }

            rr := httptest.NewRecorder()
            handler := http.HandlerFunc(LoginHandler)
            handler.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }

            if !bytes.Contains(rr.Body.Bytes(), []byte(tt.expectedBody)) {
                t.Errorf("handler returned unexpected body: got %v want %v",
                    rr.Body.String(), tt.expectedBody)
            }
        })
    }
}