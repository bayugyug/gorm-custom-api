package tools

import (
	"os"

	"github.com/bayugyug/gorm-custom-api/configs"
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

// EmptyDBTst free the test db
func (h DevelTester) EmptyDBTst(storage *gorm.DB) {
	// Select all records from a model and delete all
	// storage.Model(&models.Building{}).Delete(&models.Building{})
	storage.Exec(`
	SET FOREIGN_KEY_CHECKS = 0; 
	TRUNCATE table buildings; 
	SET FOREIGN_KEY_CHECKS = 1;
	`)
}
