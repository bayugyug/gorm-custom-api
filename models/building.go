package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	// ErrMissingRequiredParameters reqd parameter missing
	ErrMissingRequiredParameters = errors.New("missing required parameter")
	// ErrRecordsNotFound list is empty
	ErrRecordsNotFound = errors.New("record(s) not found")
	// ErrRecordNotFound data not exiss
	ErrRecordNotFound = errors.New("record not found")
	// ErrRecordMismatch generated hashkey by name is a mismatch
	ErrRecordMismatch = errors.New("record id/name mismatch")
	// ErrRecordExists data already exiss
	ErrRecordExists = errors.New("record exists")
	// ErrDBTransaction internal storage error
	ErrDBTransaction = errors.New("db storage failed")
)

// Building table buildings
type Building struct {
	ID             int64           `gorm:"primary_key" json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Address        string          `json:"address,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty"`
	BuildingFloors []BuildingFloor `gorm:"ForeignKey:BuildingID" json:"floors,omitempty"`
}

// TableName table name to be buildings
func (Building) TableName() string {
	return "buildings"
}

// BuildingFloor table building_floors
type BuildingFloor struct {
	ID         int64     `gorm:"primary_key" json:"-"`
	Floor      string    `json:"floor,omitempty"`
	BuildingID int64     `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// TableName table name to be buildings
func (BuildingFloor) TableName() string {
	return "building_floors"
}

// NewBuildingData new instance
func NewBuildingData() *Building {
	return &Building{}
}

// HashKey convert to md5 hash
func (q Building) HashKey(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GET

// BuildingGetParams get parameter
type BuildingGetParams struct {
	ID int64 `json:"id"`
}

// NewBuildingGetOne new instance with parameter
func NewBuildingGetOne(id int64) *BuildingGetParams {
	return &BuildingGetParams{ID: id}
}

// Get query from the db base on id
func (p *BuildingGetParams) Get(gDB *gorm.DB) (*Building, error) {
	// Get 1 record by id
	var building Building
	gDB.Preload("BuildingFloors").Find(&building, p.ID)
	if building.ID <= 0 {
		//not found
		return nil, ErrRecordMismatch
	}
	//yes found
	return &building, nil
}

// GetAll query all from the db
func (p *BuildingGetParams) GetAll(gDB *gorm.DB) ([]Building, error) {
	// Get all records
	var buildings []Building
	gDB.
		Preload("BuildingFloors").
		Find(&buildings)

	//empty
	if len(buildings) <= 0 {
		return buildings, ErrRecordsNotFound
	}

	return buildings, nil
}

// CREATE

// BuildingCreateParams create parameter
type BuildingCreateParams struct {
	Name    *string  `json:"name"`
	Address string   `json:"address"`
	Floors  []string `json:"floors"`
}

// NewBuildingCreate new creator
func NewBuildingCreate() *BuildingCreateParams {
	return &BuildingCreateParams{}
}

// Bind filter parameter
func (p *BuildingCreateParams) Bind(r *http.Request) error {
	//sanity check
	if p == nil {
		return ErrMissingRequiredParameters
	}
	p.Address = strings.TrimSpace(p.Address)
	//check
	return p.SanityCheck()
}

// SanityCheck filter required parameter
func (p *BuildingCreateParams) SanityCheck() error {
	if p.Name == nil || *p.Name == "" {
		return ErrMissingRequiredParameters
	}
	return nil
}

// Create add a row from the store
func (p *BuildingCreateParams) Create(gDB *gorm.DB) (int64, error) {
	//should not happen
	if err := p.SanityCheck(); err != nil {
		return 0, err
	}

	// Get 1 record by name
	var building Building
	gDB.
		Preload("BuildingFloors").
		Find(&building,
			"name = ?",
			p.Name)
	if building.ID > 0 {
		//not found
		return 0, ErrRecordExists
	}

	//  all floors
	var floors []BuildingFloor
	for _, f := range p.Floors {
		floors = append(floors,
			BuildingFloor{
				BuildingID: building.ID,
				Floor:      f,
			})
	}
	bdata := Building{
		Name:           *p.Name,
		Address:        p.Address,
		BuildingFloors: floors,
	}
	if err := gDB.Create(&bdata).Error; err != nil || bdata.ID <= 0 {
		//not found
		return 0, ErrDBTransaction
	}

	return bdata.ID, nil
}

// UPDATE

// BuildingUpdateParams update parameter
type BuildingUpdateParams struct {
	ID *int64 `json:"id"`
	BuildingCreateParams
}

// NewBuildingUpdate new instance
func NewBuildingUpdate() *BuildingUpdateParams {
	return &BuildingUpdateParams{}
}

// Bind filter parameter
func (p *BuildingUpdateParams) Bind(r *http.Request) error {
	//sanity check
	if p == nil {
		return ErrMissingRequiredParameters
	}
	//fmt
	p.Address = strings.TrimSpace(p.Address)
	//chk
	return p.SanityCheck()
}

// SanityCheck filter required parameter
func (p *BuildingUpdateParams) SanityCheck() error {
	if p.ID == nil || p.Name == nil ||
		*p.ID == 0 || *p.Name == "" {
		return ErrMissingRequiredParameters
	}
	return nil
}

// Update a row from the store
func (p *BuildingUpdateParams) Update(gDB *gorm.DB) error {
	//should not happen :-)
	if err := p.SanityCheck(); err != nil {
		return err
	}

	//1 record by id
	var building Building
	gDB.
		Set("gorm:query_option", "FOR UPDATE").
		Preload("BuildingFloors").
		Find(&building, *p.ID)

	if building.ID <= 0 {
		//not found
		return ErrRecordNotFound
	}
	//set building
	gDB.
		Model(&building).
		Updates(
			Building{
				Name:    *p.Name,
				Address: p.Address,
			})

	//remove the old floors
	gDB.
		Delete(
			BuildingFloor{},
			"building_id = ?",
			building.ID)

	//set all floors
	for _, floor := range p.Floors {
		fdata := BuildingFloor{
			BuildingID: building.ID,
			Floor:      floor,
		}
		gDB.
			Create(&fdata)
	}
	return nil
}

// DELETE

// BuildingDeleteParams delete parameter
type BuildingDeleteParams struct {
	ID int64 `json:"id"`
}

// NewBuildingDelete new instance
func NewBuildingDelete(pid int64) *BuildingDeleteParams {
	return &BuildingDeleteParams{ID: pid}
}

// Delete remove a row from the store base on id
func (p *BuildingDeleteParams) Delete(gDB *gorm.DB) error {
	// Get 1 record by id
	var building Building
	gDB.
		Preload("BuildingFloors").
		Find(&building, p.ID)
	if building.ID <= 0 {
		//not found
		return ErrRecordNotFound
	}

	//building
	gDB.
		Delete(&building)
	return nil
}