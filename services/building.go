package services

import (
	"github.com/bayugyug/gorm-custom-api/models"
	"github.com/bayugyug/gorm-custom-api/tools"
	"github.com/jinzhu/gorm"
)

// BuildingGetParams get parameter
type BuildingGetParams struct {
	ID int64 `json:"id"`
}

// NewBuildingGetOne new instance with parameter
func NewBuildingGetOne(id int64) *BuildingGetParams {
	return &BuildingGetParams{ID: id}
}

// BuildingCreateParams create parameter
type BuildingCreateParams struct {
	Name    *string  `json:"name"`
	Address string   `json:"address"`
	Floors  []string `json:"floors"`
}

// BuildingUpdateParams update parameter
type BuildingUpdateParams struct {
	ID *int64 `json:"id"`
	BuildingCreateParams
}

// BuildingDeleteParams delete parameter
type BuildingDeleteParams struct {
	ID int64 `json:"id"`
}

// NewBuildingCreate new creator
func NewBuildingCreate() *BuildingCreateParams {
	return &BuildingCreateParams{}
}

// NewBuildingUpdate new instance
func NewBuildingUpdate() *BuildingUpdateParams {
	return &BuildingUpdateParams{}
}

// NewBuildingDelete new instance
func NewBuildingDelete(pid int64) *BuildingDeleteParams {
	return &BuildingDeleteParams{ID: pid}
}

// BuildingService service creator
type BuildingService struct {
}

// NewBuildingService create a service
func NewBuildingService() *BuildingService {
	return new(BuildingService)
}

// Get query from the db base on id
func (s *BuildingService) Get(dbh *gorm.DB, p *BuildingGetParams) (*models.Building, error) {
	//prepare
	bdata := models.NewBuildingData()
	bdata.ID = p.ID
	return bdata.Get(dbh)
}

// GetAll query all from the db
func (s *BuildingService) GetAll(dbh *gorm.DB, paging *tools.PagingParams) (models.BuildingGetResults, error) {
	//prepare
	return models.NewBuildingData().GetAll(dbh, paging.Page, paging.Limit)
}

// Create add a row from the store
func (s *BuildingService) Create(dbh *gorm.DB, p *BuildingCreateParams) (int64, error) {
	//prepare
	var floors []models.BuildingFloor
	for _, f := range p.Floors {
		floors = append(floors,
			models.BuildingFloor{
				Floor: f,
			})
	}
	bdata := models.Building{
		Name:           *p.Name,
		Address:        p.Address,
		BuildingFloors: floors,
	}
	return bdata.Create(dbh)
}

// Update modify old data
func (s *BuildingService) Update(dbh *gorm.DB, p *BuildingUpdateParams) error {
	//prepare
	var floors []models.BuildingFloor
	for _, f := range p.Floors {
		floors = append(floors,
			models.BuildingFloor{
				BuildingID: *p.ID,
				Floor:      f,
			})
	}
	bdata := &models.Building{
		ID:             *p.ID,
		Name:           *p.Name,
		Address:        p.Address,
		BuildingFloors: floors,
	}
	return bdata.Update(dbh)
}

// Delete remove a row from the store base on id
func (s *BuildingService) Delete(dbh *gorm.DB, p *BuildingDeleteParams) error {
	//prepare
	bdata := models.NewBuildingData()
	bdata.ID = p.ID
	return bdata.Delete(dbh)
}
