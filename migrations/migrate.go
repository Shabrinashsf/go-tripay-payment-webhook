package migrations

import (
	"github.com/Shabrinashsf/go-tripay-payment-webhook/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := db.AutoMigrate(
		&entity.Product{},
		&entity.Transaction{},
	); err != nil {
		return err
	}

	return nil
}
