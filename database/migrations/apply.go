package migrations

import (
	"context"

	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"gorm.io/gorm"
)

func Apply(ctx context.Context, db *gorm.DB) error {
	ctx = context.WithValue(ctx, utils.DBContextKey, db)

	if err := Up00001Init(ctx); err != nil {
		return err
	}
	// Down00001Init(ctx)

	return nil
}
