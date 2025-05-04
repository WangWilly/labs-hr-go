package employee

import (
	"net/http"
	"time"

	"github.com/WangWilly/labs-gin/pkgs/models"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	Name    string `json:"name"    binding:"required"`
	Age     int    `json:"age"     binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone"   binding:"required"`
	Email   string `json:"email"   binding:"required"`

	Position   string  `json:"position"   binding:"required"`
	Department string  `json:"department" binding:"required"`
	Salary     float64 `json:"salary"     binding:"required"`
	StartDate  int64   `json:"start_date" binding:"required"`
}

type CreateResponse struct {
	EmployeeID int64 `json:"employee_id"`
	PositionID int64 `json:"position_id"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Create the employee info
	employeeInfo := &models.EmployeeInfo{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
	}
	if err := c.employeeInfoRepo.Create(ctx, c.db, employeeInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee info"})
		return
	}

	// Create the employee position
	employeePosition := &models.EmployeePosition{
		EmployeeID: employeeInfo.ID,
		Position:   req.Position,
		Department: req.Department,
		Salary:     req.Salary,
		StartDate:  time.Unix(req.StartDate, 0),
	}
	nowTime := c.timeModule.Now()
	if err := c.employeePositionRepo.Create(ctx, c.db, employeePosition, nowTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee position"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	ctx.JSON(http.StatusCreated, CreateResponse{
		EmployeeID: employeeInfo.ID,
		PositionID: employeePosition.ID,
	})
}
