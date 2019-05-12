package tools

import (
	"fmt"
	"math/rand"

	"github.com/icrowley/fake"
)

// Seeder generator of dummy data
type Seeder struct {
}

// Create geneate dummy building-data
func (s Seeder) Create() string {
	return fmt.Sprintf(`{ "name": "building-%s","address": "address::%s","floors": ["floor-%s","floor-%s"] }`,
		fake.DigitsN(12),
		fake.DigitsN(15),
		fake.DigitsN(5),
		fake.DigitsN(5),
	)
}

// CreateWithName geneate dummy building-data and with name
func (s Seeder) CreateWithName(name string) string {
	return fmt.Sprintf(`{ "name": "%s","address": "address::%s","floors": ["floor-%s","floor-%s"] }`,
		name,
		fake.DigitsN(15),
		fake.DigitsN(5),
		fake.DigitsN(5),
	)
}

// Update geneate dummy building-data with an ID from parameter
func (s Seeder) Update(pid int64, name string) string {
	return fmt.Sprintf(`{ "id": %d, "name": "%s","address": "address::%s","floors": ["floor-%s","floor-%s"] }`,
		pid,
		name,
		fake.DigitsN(15),
		fake.DigitsN(5),
		fake.DigitsN(5),
	)
}

// CreateWithEmptyName  geneate dummy building-data with empty name
func (s Seeder) CreateWithEmptyName() string {
	return fmt.Sprintf(`{ "name": "","address": "address::%s","floors": ["floor-%s","floor-%s"] }`,
		fake.DigitsN(15),
		fake.DigitsN(5),
		fake.DigitsN(5),
	)
}

// CreateFloors generate random floor list
func (s Seeder) CreateFloors() []string {
	var floors []string
	t := rand.Intn(10) + 1
	for i := 1; i <= t; i++ {
		floors = append(floors, fmt.Sprintf("floor-%s", fake.DigitsN(15)))
	}
	return floors
}

// CreateMin geneate dummy building-data and with only name
func (s Seeder) CreateMin() string {
	return fmt.Sprintf(`{ "name": "anonymous building %s" }`,
		fake.DigitsN(15),
	)
}
