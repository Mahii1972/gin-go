package handlers

import (
	"GIN/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmployeeHandler struct {
	queries *db.Queries
}

func NewEmployeeHandler(queries *db.Queries) *EmployeeHandler {
	return &EmployeeHandler{
		queries: queries,
	}
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.queries.ListEmployees(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, employees)
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	employee, err := h.queries.GetEmployee(c, int32(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(404, gin.H{"error": "Employee not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, employee)
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	employee, err := h.queries.CreateEmployee(c, db.CreateEmployeeParams(req))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	employee, err := h.queries.UpdateEmployee(c, db.UpdateEmployeeParams{
		ID: int32(id), FirstName: req.FirstName, LastName: req.LastName,
		Email: req.Email, HireDate: req.HireDate, Salary: req.Salary,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, employee)
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.queries.DeleteEmployee(c, int32(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Employee deleted successfully"})
}

// Move the EmployeeRequest struct here
type EmployeeRequest struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	HireDate  pgtype.Date    `json:"hire_date"`
	Salary    pgtype.Numeric `json:"salary"`
}
