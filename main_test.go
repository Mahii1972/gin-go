package main

import (
	"GIN/db"
	"GIN/handlers"
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	testRouter  *gin.Engine
	testQueries *db.Queries
)

func TestMain(m *testing.M) {
	// Setup phase
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	testQueries = db.New(conn)
	handler := handlers.NewEmployeeHandler(testQueries)

	testRouter = gin.Default()
	testRouter.GET("/employees", handler.ListEmployees)
	testRouter.GET("/employee/:id", handler.GetEmployee)
	testRouter.POST("/employee", handler.CreateEmployee)
	testRouter.PUT("/employee/:id", handler.UpdateEmployee)
	testRouter.DELETE("/employee/:id", handler.DeleteEmployee)

	// Run tests
	code := m.Run()

	// Exit with the test status code
	os.Exit(code)
}

func TestGetEmployee(t *testing.T) {
	tests := []struct {
		name           string
		employeeID     int
		expectedCode   int
		expectedError  bool
		expectedFields []string
	}{
		{
			name:           "Valid Employee ID",
			employeeID:     5,
			expectedCode:   http.StatusOK,
			expectedError:  false,
			expectedFields: []string{"ID", "FirstName", "LastName", "Email", "HireDate", "Salary"},
		},
		{
			name:           "Non-existent Employee ID",
			employeeID:     999999,
			expectedCode:   http.StatusNotFound,
			expectedError:  true,
			expectedFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/employee/"+strconv.Itoa(tt.employeeID), nil)
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Contains(t, response, "error")
			} else {
				assert.NotContains(t, response, "error")
				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field)
				}
			}
		})
	}
}

func TestListEmployees(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/employees", nil)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Greater(t, len(response), 0)
}

func TestCreateEmployee(t *testing.T) {
	newEmployee := handlers.EmployeeRequest{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice.smith@example.com",
		HireDate:  pgtype.Date{Time: time.Now(), Valid: true},
		Salary:    pgtype.Numeric{Int: big.NewInt(50000), Valid: true},
	}

	body, _ := json.Marshal(newEmployee)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/employee", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, newEmployee.FirstName, response["FirstName"])
	assert.Equal(t, newEmployee.LastName, response["LastName"])
	assert.Equal(t, newEmployee.Email, response["Email"])
}

func TestUpdateEmployee(t *testing.T) {
	updateEmployee := handlers.EmployeeRequest{
		FirstName: "Alice",
		LastName:  "Johnson",
		Email:     "alice.johnson@example.com",
		HireDate:  pgtype.Date{Time: time.Now(), Valid: true},
		Salary:    pgtype.Numeric{Int: big.NewInt(55000), Valid: true},
	}

	body, _ := json.Marshal(updateEmployee)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/employee/5", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateEmployee.FirstName, response["FirstName"])
	assert.Equal(t, updateEmployee.LastName, response["LastName"])
	assert.Equal(t, updateEmployee.Email, response["Email"])
}

func TestDeleteEmployee(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/employee/4", nil)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Employee deleted successfully", response["message"])
}
