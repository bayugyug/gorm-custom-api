package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bayugyug/gorm-custom-api/drivers"

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

//BuildingCreator list of all building funcs
type BuildingCreator interface {
	Get(dbh *gorm.DB) (*Building, error)
	GetAll(dbh *gorm.DB) ([]Building, error)
	Create(dbh *gorm.DB) (int64, error)
	Update(dbh *gorm.DB) error
	Delete(dbh *gorm.DB) error
}

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

// BuildingGetResults table buildings query list
type BuildingGetResults struct {
	Items []Building `json:"items,omitempty"`
	Page  int        `json:"page,omitempty"`
	Limit int        `json:"limit,omitempty"`
	Total int        `json:"total,omitempty"`
}

// NewBuildingData new instance
func NewBuildingData() *Building {
	return &Building{}
}

// HashKey convert to md5 hash
func (q *Building) HashKey(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// Get query from the db base on id
func (q *Building) Get(dbh *gorm.DB) (*Building, error) {
	// Get 1 record by id
	var building Building
	dbh.
		Preload("BuildingFloors").
		Find(&building, q.ID)
	if building.ID <= 0 {
		//not found
		return nil, ErrRecordMismatch
	}
	//yes found
	return &building, nil
}

// GetAll query all from the db
func (q *Building) GetAll(dbh *gorm.DB, page, limit int) (BuildingGetResults, error) {
	// Get all records
	var buildings []Building
	total := 0
	dbh.
		Preload("BuildingFloors").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&buildings).
		Offset(0).
		Model(&Building{}).
		Limit(-1).
		Count(&total)

	log.Println("PAGE: ", page, "LIMIT: ", limit, ", TOTAL: ", total)
	//empty
	if len(buildings) <= 0 {
		return BuildingGetResults{}, ErrRecordsNotFound
	}

	return BuildingGetResults{
		Items: buildings,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

// Create add a row from the store
func (q *Building) Create(dbh *gorm.DB) (int64, error) {
	// Get 1 record by name
	var err error
	var building Building
	//exists
	dbh.
		Preload("BuildingFloors").
		Find(&building,
			"name = ?",
			q.Name)
	if building.ID > 0 {
		//oops found
		return 0, ErrRecordExists
	}
	//sync run
	err = drivers.SyncRunTx(dbh, func(transaction *gorm.DB) error {
		return transaction.Create(q).Error
	})

	if err != nil || q.ID <= 0 {
		log.Println("Failed::Create", err)
		//not found
		return 0, ErrDBTransaction
	}
	return q.ID, nil
}

// Update modify old data
func (q *Building) Update(dbh *gorm.DB) error {
	var err error
	//1 record by id
	var building Building
	dbh.
		Set("gorm:query_option", "FOR UPDATE").
		Preload("BuildingFloors").
		Find(&building, q.ID)

	//not found
	if building.ID <= 0 {
		return ErrRecordNotFound
	}

	//sync run
	err = drivers.SyncRunTx(dbh, func(transaction *gorm.DB) error {
		terr := transaction.
			Model(&building).
			Updates(
				Building{
					Name:    q.Name,
					Address: q.Address,
				}).Error

		if terr != nil {
			return terr
		}

		//remove the old floors
		terr = transaction.
			Delete(
				BuildingFloor{},
				"building_id = ?",
				building.ID).Error

		if terr != nil {
			//not found
			return ErrDBTransaction
		}
		return nil
	})

	if err != nil {
		log.Println("Failed::Update", err)
		return ErrDBTransaction
	}

	//set all floors
	for _, f := range q.BuildingFloors {
		fdata := BuildingFloor{
			BuildingID: building.ID,
			Floor:      f.Floor,
		}

		err = drivers.SyncRunTx(dbh, func(transaction *gorm.DB) error {
			return transaction.Create(&fdata).Error
		})

		if err != nil {
			log.Println("Failed::Update", err)
			return ErrDBTransaction
		}
	}
	return nil
}

// Delete remove a row from the store base on id
func (q *Building) Delete(dbh *gorm.DB) error {
	var err error
	// Get 1 record by id
	var building Building
	dbh.
		Preload("BuildingFloors").
		Find(&building, q.ID)

	//not found
	if building.ID <= 0 {
		return ErrRecordNotFound
	}

	//building
	err = drivers.SyncRunTx(dbh, func(transaction *gorm.DB) error {
		return transaction.
			Delete(&building).
			Error
	})

	if err != nil {
		log.Println("Failed::Update", err)
		return ErrDBTransaction
	}
	return nil
}

// CALLBACKS

// BeforeSave callback before save
func (q *Building) BeforeSave() (err error) {
	return
}

// AfterSave callback after save
func (q *Building) AfterSave(dbh *gorm.DB) (err error) {
	return
}

// BeforeCreate callback before create
func (q *Building) BeforeCreate() (err error) {
	return
}

// AfterCreate callback after create
func (q *Building) AfterCreate(dbh *gorm.DB) (err error) {
	return
}

// BeforeUpdate callback before update
func (q *Building) BeforeUpdate() (err error) {
	return
}

// AfterUpdate callback after update
func (q *Building) AfterUpdate(dbh *gorm.DB) (err error) {
	return
}

// BeforeDelete callback before remove
func (q *Building) BeforeDelete() (err error) {
	return
}

// AfterDelete callback after remove
func (q *Building) AfterDelete(dbh *gorm.DB) (err error) {
	return
}
