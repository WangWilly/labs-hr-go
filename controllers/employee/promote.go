package employee

import (
	"strconv"
	"time"

	"github.com/WangWilly/labs-gin/pkgs/models"
	"github.com/WangWilly/labs-gin/pkgs/utils"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type PromoteRequest struct {
	Position   string  `json:"position"   binding:"required"`
	Department string  `json:"department" binding:"required"`
	Salary     float64 `json:"salary"     binding:"required"`
	StartDate  int64   `json:"start_date" binding:"required"`
}

type PromoteResponse struct {
	PositionID int64  `json:"position_id"`
	StartDate  string `json:"start_date"`
}

func (c *Controller) Promote(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	// Convert id to int64
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var req PromoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	employeePosition := &models.EmployeePosition{
		EmployeeID: employeeID,
		Position:   req.Position,
		Department: req.Department,
		Salary:     req.Salary,
		StartDate:  time.Unix(req.StartDate, 0),
	}
	nowTime := c.timeModule.Now()
	if err := c.employeePositionRepo.Create(ctx, c.db, employeePosition, nowTime); err != nil {
		ctx.JSON(500, gin.H{"error": "failed to create employee position"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	response := PromoteResponse{
		PositionID: employeePosition.ID,
		StartDate:  utils.FormatedTime(employeePosition.StartDate),
	}
	ctx.JSON(200, response)
}
