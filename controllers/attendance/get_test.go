package attendance

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WangWilly/labs-gin/pkgs/models"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	testInit(t, func(s *testSuite) {
		Convey("Given an employee exists in the system", t, func() {
			// Setup test data
			employeeID := int64(123)
			positionID := int64(456)
			attendanceID := int64(789)
			nowTime := time.Date(2023, 6, 15, 9, 0, 0, 0, time.UTC)

			employeePosition := &models.EmployeePosition{
				ID:         positionID,
				EmployeeID: employeeID,
				Position:   "Software Engineer",
				Department: "Engineering",
				Salary:     90000.00,
				StartDate:  nowTime.Add(-30 * 24 * time.Hour), // 30 days ago
			}

			attendance := &models.EmployeeAttendance{
				ID:         attendanceID,
				EmployeeID: employeeID,
				PositionID: positionID,
				ClockIn:    nowTime.Add(-8 * time.Hour),
				ClockOut:   nowTime,
			}

			Convey("When retrieving the employee's attendance record", func() {
				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(attendance, nil)

				// Expected response
				expectedResponse := GetResponse{
					AttendanceID: attendanceID,
					PositionID:   positionID,
					ClockInTime:  "2023-06-15 01:00:00",
					ClockOutTime: "2023-06-15 09:00:00",
				}

				// Make the request and verify response
				var actualResponse GetResponse
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&actualResponse,
					http.StatusOK,
				)

				Convey("Then the response should contain the attendance details", func() {
					So(actualResponse.AttendanceID, ShouldEqual, expectedResponse.AttendanceID)
					So(actualResponse.PositionID, ShouldEqual, expectedResponse.PositionID)
					So(actualResponse.ClockInTime, ShouldEqual, expectedResponse.ClockInTime)
					So(actualResponse.ClockOutTime, ShouldEqual, expectedResponse.ClockOutTime)
				})
			})

			Convey("When retrieving with invalid employee ID format", func() {
				// Make the request with invalid ID and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/invalid-id",
					nil,
					&errorResponse,
					http.StatusBadRequest,
				)

				Convey("Then the response should indicate invalid employee ID", func() {
					So(errorResponse["error"], ShouldEqual, "invalid employee id")
				})
			})

			Convey("When employee position is not found", func() {
				// Set up expectations for failure case
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(nil, nil)

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate position not found", func() {
					So(errorResponse["error"], ShouldEqual, "employee position not found")
				})
			})

			Convey("When getting position fails with database error", func() {
				// Set up expectations for database error
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get employee position")
				})
			})

			Convey("When attendance record is not found", func() {
				// Set up expectations
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, nil)

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusNotFound,
				)

				Convey("Then the response should indicate attendance not found", func() {
					So(errorResponse["error"], ShouldEqual, "attendance not found")
				})
			})

			Convey("When getting attendance fails with database error", func() {
				// Set up expectations for database error
				s.timeModule.EXPECT().Now().Return(nowTime)

				s.employeePositionRepo.EXPECT().
					GetCurrentByEmployeeID(gomock.Any(), s.db, employeeID, nowTime).
					Return(employeePosition, nil)

				s.employeeAttendanceRepo.EXPECT().
					Last(gomock.Any(), s.db, employeeID).
					Return(nil, errors.New("database error"))

				// Make the request and verify error response
				var errorResponse map[string]string
				s.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/attendance/123",
					nil,
					&errorResponse,
					http.StatusInternalServerError,
				)

				Convey("Then the response should indicate a server error", func() {
					So(errorResponse["error"], ShouldEqual, "failed to get last attendance")
				})
			})
		})
	})
}
