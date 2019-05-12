package tools

import (
	"crypto/md5"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/bayugyug/gorm-custom-api/configs"
	"github.com/bayugyug/gorm-custom-api/models"
	"github.com/google/uuid"
)

// Helper utility
type Helper struct {
}

// RemoveIntDuplicates de-duplicate int list
func (h Helper) RemoveIntDuplicates(elements []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

// RemoveStrDuplicates de-duplicate string list
func (h Helper) RemoveStrDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

// FormatSliceToIntMap convert slice of int to map
func (h Helper) FormatSliceToIntMap(all []int) map[int]int {
	bmap := make(map[int]int)
	for _, bv := range all {
		bmap[bv] = bv
	}
	return bmap
}

// UUID random uuid
func (h Helper) UUID() string {
	return uuid.New().String() + `-` + time.Now().Format("20060102-150405")
}

// HashMD5 generate md5 string
func (h Helper) HashMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// ConfigTst set up the config for testing from envt
func (h Helper) ConfigTst() *configs.ParameterConfig {
	var s = `{"port":"8989","dsn":"root:@tcp(127.0.0.1:3306)/gorm_cust_api_test"}`
	cfg := &configs.APISettings{}
	if os.Getenv("GORM_CUST_API_CONFIG") != "" {
		s = os.Getenv("GORM_CUST_API_CONFIG")
	}
	return cfg.FormatParameterConfig(s)
}

// EmptyDBTst free the test db
func (h Helper) EmptyDBTst(storage *gorm.DB) {
	// Select all records from a model and delete all
	storage.Model(&models.Building{}).Delete(&models.Building{})
}
