package attendance

import (
	"net/http"
	"strconv"

	"github.com/WangWilly/labs-gin/pkgs/utils"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type GetResponse struct {
	AttendanceID int64  `json:"attendance_id"`
	PositionID   int64  `json:"position_id"`
	ClockInTime  string `json:"clock_in_time"`
	ClockOutTime string `json:"clock_out_time"`
}

func (c *Controller) Get(ctx *gin.Context) {
	employeeID, ok := ctx.Params.Get("employee_id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	employeeIDInt, err := strconv.ParseInt(employeeID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Get the current position of the employee
	employeePosition, err := c.employeePositionRepo.GetCurrentByEmployeeID(ctx, c.db, employeeIDInt, c.timeModule.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get employee position"})
		return
	}
	if employeePosition == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "employee position not found"})
		return
	}

	// Get the last attendance record of the employee
	currAttendance, err := c.employeeAttendanceRepo.Last(ctx, c.db, employeeIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get last attendance"})
		return
	}
	if currAttendance == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}

	clockOutTime := ""
	if currAttendance.ClockIn != currAttendance.ClockOut {
		// If the clock-in and clock-out times are the same, it means the employee has not clocked out yet
		clockOutTime = utils.FormatedTime(currAttendance.ClockOut)
	}

	////////////////////////////////////////////////////////////////////////////
	// Return the attendance record
	ctx.JSON(http.StatusOK, GetResponse{
		AttendanceID: currAttendance.ID,
		PositionID:   currAttendance.PositionID,
		ClockInTime:  utils.FormatedTime(currAttendance.ClockIn),
		ClockOutTime: clockOutTime,
	})
}
