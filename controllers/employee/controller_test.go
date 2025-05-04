package employee

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WangWilly/labs-hr-go/pkgs/testutils"
	"github.com/brianvoe/gofakeit/v6"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	db     *gorm.DB
	mockDB sqlmock.Sqlmock

	timeModule           *MockTimeModule
	employeeInfoRepo     *MockEmployeeInfoRepo
	employeePositionRepo *MockEmployeePositionRepo

	controller *Controller
	testServer testutils.TestHttpServer
	faker      *gofakeit.Faker
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mockDB, _ := sqlmock.New()
	gormDB, _ := gorm.Open(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{SkipDefaultTransaction: true},
	)

	timeModule := NewMockTimeModule(ctrl)
	employeeInfoRepo := NewMockEmployeeInfoRepo(ctrl)
	employeePositionRepo := NewMockEmployeePositionRepo(ctrl)

	cfg := Config{}
	controller := NewController(
		cfg,
		gormDB,
		timeModule,
		employeeInfoRepo,
		employeePositionRepo,
	)
	testServer := testutils.NewTestHttpServer(controller)
	faker := gofakeit.New(0)
	suite := &testSuite{
		db:                   gormDB,
		mockDB:               mockDB,
		timeModule:           timeModule,
		employeeInfoRepo:     employeeInfoRepo,
		employeePositionRepo: employeePositionRepo,
		controller:           controller,
		testServer:           testServer,
		faker:                faker,
	}

	test(suite)
}
