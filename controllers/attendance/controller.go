package attendance

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
}

type Controller struct {
	cfg Config
	db  *gorm.DB

	timeModule             TimeModule
	employeePositionRepo   EmployeePositionRepo
	employeeAttendanceRepo EmployeeAttendanceRepo
}

func NewController(
	cfg Config,
	db *gorm.DB,
	timeModule TimeModule,
	employeePositionRepo EmployeePositionRepo,
	employeeAttendanceRepo EmployeeAttendanceRepo,
) *Controller {
	return &Controller{
		cfg:                    cfg,
		db:                     db,
		timeModule:             timeModule,
		employeePositionRepo:   employeePositionRepo,
		employeeAttendanceRepo: employeeAttendanceRepo,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	////////////////////////////////////////////////////////////////////////////
	// attendance management
	r.POST("/attendance", c.Create)
	r.GET("/attendance/:employee_id", c.Get)
}
