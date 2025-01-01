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
func setupTestDB() {
    var err error
    database.Connection, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    database.Connection.AutoMigrate(&models.User{}, &models.Token{})
}

func TestIntegration_LoginLogout(t *testing.T) {
    setupTestDB()

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

    // Test Login
    loginRequest := LoginRequest{
        Username: "testuser",
        Password: "testpassword",
    }
    loginBody, _ := json.Marshal(loginRequest)
    loginReq, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
    if err != nil {
        t.Fatal(err)
    }
    loginRR := httptest.NewRecorder()
    loginHandler := http.HandlerFunc(LoginHandler)
    loginHandler.ServeHTTP(loginRR, loginReq)

    if status := loginRR.Code; status != http.StatusOK {
        t.Errorf("LoginHandler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    var loginResponse map[string]interface{}
    json.NewDecoder(loginRR.Body).Decode(&loginResponse)
    token, ok := loginResponse["token"].(string)
    if !ok {
        t.Fatal("LoginHandler did not return a token")
    }

    // Test Logout
    logoutReq, err := http.NewRequest("POST", "/logout", nil)
    if err != nil {
        t.Fatal(err)
    }
    logoutReq.Header.Set("Authorization", "Bearer "+token)
    logoutRR := httptest.NewRecorder()
    logoutHandler := http.HandlerFunc(LogoutHandler)
    logoutHandler.ServeHTTP(logoutRR, logoutReq)

    if status := logoutRR.Code; status != http.StatusOK {
        t.Errorf("LogoutHandler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    if !bytes.Contains(logoutRR.Body.Bytes(), []byte(`"status":"success"`)) {
        t.Errorf("LogoutHandler returned unexpected body: got %v want %v",
            logoutRR.Body.String(), `"status":"success"`)
    }
}