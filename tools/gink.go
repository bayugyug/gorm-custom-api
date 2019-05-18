package tools

import (
	"os"

	"github.com/bayugyug/gorm-custom-api/configs"
	"github.com/bayugyug/gorm-custom-api/models"
	"github.com/jinzhu/gorm"
)

// DevelTester for the setup of ginkgo test
type DevelTester struct {
}

// Config set up the config for testing from envt
func (h DevelTester) Config() *configs.ParameterConfig {
	//make sure TEST_GORM_CUST_API_CONFIG is defined from envt
	cfg := &configs.APISettings{}
	return cfg.FormatParameterConfig(os.Getenv("TEST_GORM_CUST_API_CONFIG"))
}

// EmptyDBTst free the test db
func (h DevelTester) EmptyDBTst(storage *gorm.DB) {
	// Select all records from a model and delete all
	storage.Model(&models.Building{}).Delete(&models.Building{})
	// storage.Exec("DELETE FROM buildings;")
}
