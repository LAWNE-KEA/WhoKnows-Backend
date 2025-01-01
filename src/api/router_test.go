// package api

// file looks correct, needs minor adjustments and it should work

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"whoKnows/database"
// 	"whoKnows/models"
// 	"whoKnows/security"
// 	"whoKnows/api/handlers" // Import the handlers package
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func init() {
// 	// Initialize the in-memory database for tests
// 	var err error
// 	database.Connection, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect to the database")
// 	}
// 	database.Connection.AutoMigrate(&models.User{})
// }

// func TestLoginHandler(t *testing.T) {
// 	// Create a mock user with hashed password
// 	hashedPassword, err := security.HashPassword("testpassword")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	user := models.User{
// 		Username: "testuser",
// 		Password: hashedPassword,
// 	}
// 	database.Connection.Create(&user)

// 	tests := []struct {
// 		name           string
// 		requestBody    handlers.LoginRequest // Use handlers.LoginRequest
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Valid login",
// 			requestBody: handlers.LoginRequest{ // Use handlers.LoginRequest
// 				Username: "testuser",
// 				Password: "testpassword",
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `"status":"success"`, // Expecting success status
// 		},
// 		{
// 			name: "Invalid login - missing username",
// 			requestBody: handlers.LoginRequest{ // Use handlers.LoginRequest
// 				Password: "testpassword",
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `"Username and password are required"`,
// 		},
// 		{
// 			name: "Invalid login - wrong password",
// 			requestBody: handlers.LoginRequest{ // Use handlers.LoginRequest
// 				Username: "testuser",
// 				Password: "wrongpassword",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   `"Invalid username or password"`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Marshal the request body
// 			body, _ := json.Marshal(tt.requestBody)
// 			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			// Create a response recorder to capture the response
// 			rr := httptest.NewRecorder()
// 			handler := http.HandlerFunc(handlers.LoginHandler) // Correct reference to LoginHandler
// 			handler.ServeHTTP(rr, req)

// 			// Check the status code
// 			if status := rr.Code; status != tt.expectedStatus {
// 				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
// 			}

// 			// Check the response body for error or success messages
// 			if !bytes.Contains(rr.Body.Bytes(), []byte(tt.expectedBody)) {
// 				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
// 			}

// 			// Additional check for successful login
// 			if tt.name == "Valid login" {
// 				var response map[string]interface{}
// 				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
// 					t.Fatal(err)
// 				}

// 				// Check if JWT token exists in the response
// 				if token, ok := response["token"]; !ok || token == "" {
// 					t.Error("Expected token in response, but got none")
// 				}
// 			}
// 		})
// 	}
// }
