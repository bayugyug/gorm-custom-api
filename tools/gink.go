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
	var s = `{"port":"8989","dsn":"root:@tcp(127.0.0.1:3306)/gorm_cust_api_test"}`
	cfg := &configs.APISettings{}
	if os.Getenv("TEST_GORM_CUST_API_CONFIG") != "" {
		s = os.Getenv("TEST_GORM_CUST_API_CONFIG")
	}
	return cfg.FormatParameterConfig(s)
}

// Empty free the test db
func (h DevelTester) Empty(storage *gorm.DB) {
	// Select all records from a model and delete all
	storage.Model(&models.Building{}).Delete(&models.Building{})
}

// EmptyDBTst free the test db
func (h Helper) EmptyDBTst(storage *gorm.DB) {
	// Select all records from a model and delete all
	storage.Model(&models.Building{}).Delete(&models.Building{})
}
