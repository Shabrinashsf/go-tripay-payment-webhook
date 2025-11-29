package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/entity"
	"gorm.io/gorm"
)

func ListProductSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/product.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listProduct []entity.Product
	if err := json.Unmarshal(jsonData, &listProduct); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Product{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Product{}); err != nil {
			return err
		}
	}

	for _, product := range listProduct {
		if err := db.Save(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
