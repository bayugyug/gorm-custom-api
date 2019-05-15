package handler

import (
	"net/http"
	"strings"

	"github.com/bayugyug/gorm-custom-api/models"
)

// BuildingFormCreateParams create parameter
type BuildingFormCreateParams struct {
	Name    *string  `json:"name"`
	Address string   `json:"address"`
	Floors  []string `json:"floors"`
}

// NewBuildingFormCreateParams new creator
func NewBuildingFormCreateParams() *BuildingFormCreateParams {
	return &BuildingFormCreateParams{}
}

// Bind filter parameter
func (p *BuildingFormCreateParams) Bind(r *http.Request) error {
	//sanity check
	if p == nil {
		return models.ErrMissingRequiredParameters
	}
	p.Address = strings.TrimSpace(p.Address)
	//check
	return p.SanityCheck()
}

// SanityCheck filter required parameter
func (p *BuildingFormCreateParams) SanityCheck() error {
	if p.Name == nil || *p.Name == "" {
		return models.ErrMissingRequiredParameters
	}
	return nil
}

// BuildingFormUpdateParams update parameter
type BuildingFormUpdateParams struct {
	ID *int64 `json:"id"`
	BuildingFormCreateParams
}

// NewBuildingFormUpdateParams new instance
func NewBuildingFormUpdateParams() *BuildingFormUpdateParams {
	return &BuildingFormUpdateParams{}
}

// Bind filter parameter
func (p *BuildingFormUpdateParams) Bind(r *http.Request) error {
	//sanity check
	if p == nil {
		return models.ErrMissingRequiredParameters
	}
	//fmt
	p.Address = strings.TrimSpace(p.Address)
	//chk
	return p.SanityCheck()
}

// SanityCheck filter required parameter
func (p *BuildingFormUpdateParams) SanityCheck() error {
	if p.ID == nil || p.Name == nil ||
		*p.ID == 0 || *p.Name == "" {
		return models.ErrMissingRequiredParameters
	}
	return nil
}

// BuildingFormGetParams get parameter
type BuildingFormGetParams struct {
	ID int64 `json:"id"`
}

// NewBuildingFormGetParams new instance with parameter
func NewBuildingFormGetParams(id int64) *BuildingFormGetParams {
	return &BuildingFormGetParams{ID: id}
}

// BuildingFormDeleteParams delete parameter
type BuildingFormDeleteParams struct {
	ID int64 `json:"id"`
}

// NewBuildingFormDeleteParams new instance
func NewBuildingFormDeleteParams(pid int64) *BuildingFormDeleteParams {
	return &BuildingFormDeleteParams{ID: pid}
}
