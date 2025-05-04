package attendance

import (
	"net/http"

	"github.com/WangWilly/labs-gin/pkgs/utils"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	EmployeeID int64 `json:"employee_id" binding:"required"`
}

type CreateResponse struct {
	AttendanceID int64  `json:"attendance_id"`
	PositionID   int64  `json:"position_id"`
	ClockInTime  string `json:"clock_in_time"`
	ClockOutTime string `json:"clock_out_time"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Get the current position of the employee
	employeePosition, err := c.employeePositionRepo.GetCurrentByEmployeeID(ctx, c.db, req.EmployeeID, c.timeModule.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get employee position"})
		return
	}
	if employeePosition == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "employee position not found"})
		return
	}

	currAttendance, err := c.employeeAttendanceRepo.Last(ctx, c.db, req.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get last attendance"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	if currAttendance == nil {
		// Create a new attendance record for clock-in
		newAttendance, err := c.employeeAttendanceRepo.CreateForClockIn(ctx, c.db, req.EmployeeID, employeePosition.ID, c.timeModule.Now())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance"})
			return
		}
		ctx.JSON(http.StatusCreated, CreateResponse{
			AttendanceID: newAttendance.ID,
			PositionID:   newAttendance.PositionID,
			ClockInTime:  utils.FormatedTime(newAttendance.ClockIn),
			ClockOutTime: "",
		})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	if currAttendance.ClockIn == currAttendance.ClockOut {
		currAttendance, err = c.employeeAttendanceRepo.UpdateForClockOut(ctx, c.db, currAttendance.ID, c.timeModule.Now())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update attendance"})
			return
		}
		ctx.JSON(http.StatusCreated, CreateResponse{
			AttendanceID: currAttendance.ID,
			PositionID:   currAttendance.PositionID,
			ClockInTime:  utils.FormatedTime(currAttendance.ClockIn),
			ClockOutTime: utils.FormatedTime(currAttendance.ClockOut),
		})
		return
	}

	newAttendance, err := c.employeeAttendanceRepo.CreateForClockIn(ctx, c.db, req.EmployeeID, employeePosition.ID, c.timeModule.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendance"})
		return
	}
	ctx.JSON(http.StatusCreated, CreateResponse{
		AttendanceID: newAttendance.ID,
		PositionID:   newAttendance.PositionID,
		ClockInTime:  utils.FormatedTime(newAttendance.ClockIn),
		ClockOutTime: "",
	})
}
