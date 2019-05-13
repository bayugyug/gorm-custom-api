package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bayugyug/gorm-custom-api/configs"
	"github.com/bayugyug/gorm-custom-api/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
)

// BuildingEndpoints the end-points-url mapping
type BuildingEndpoints interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// Response is the reply object
type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Total  int         `json:"total,omitempty"`
}

// Building the api handler
type Building struct {
	Storage *gorm.DB
}

// NewBuilding new instance
func NewBuilding(conn *gorm.DB) *Building {
	return &Building{
		Storage: conn,
	}
}

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

// Welcome index page
func (b *Building) Welcome(w http.ResponseWriter, r *http.Request) {
	//good
	render.JSON(w, r, Response{
		Status: "Welcome!",
	})
}

// Create save a row in store
func (b *Building) Create(w http.ResponseWriter, r *http.Request) {
	fdata := NewBuildingFormCreateParams()
	//sanity check
	if err := render.Bind(r, fdata); err != nil {
		log.Println("CREATE", fdata, err)
		//400
		b.ReplyErrContent(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	data := models.NewBuildingCreate()
	data.Name = fdata.Name
	data.Address = fdata.Address
	data.Floors = fdata.Floors

	pid, err := data.Create(b.Storage)
	//chk
	if err != nil {
		log.Println("CREATE", err)
		switch err {
		case models.ErrRecordExists:
			//409
			b.ReplyErrContent(w, r, http.StatusConflict, err.Error())
		case models.ErrMissingRequiredParameters:
			//400
			b.ReplyErrContent(w, r, http.StatusBadRequest, err.Error())
		default:
			//500
			b.ReplyErrContent(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	//good
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, Response{
		Status: "success",
		Result: pid,
	})
}

// Update update row in store
func (b *Building) Update(w http.ResponseWriter, r *http.Request) {
	fdata := NewBuildingFormUpdateParams()
	//sanity check
	if err := render.Bind(r, fdata); err != nil {
		log.Println("UPDATE", fdata, err)
		//400
		b.ReplyErrContent(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	data := models.NewBuildingUpdate()
	data.ID = fdata.ID
	data.Name = fdata.Name
	data.Address = fdata.Address
	data.Floors = fdata.Floors
	//check
	if err := data.Update(b.Storage); err != nil {
		log.Println("UPDATE", err)
		switch err {
		case models.ErrRecordMismatch:
			//409
			b.ReplyErrContent(w, r, http.StatusConflict, err.Error())
		case models.ErrMissingRequiredParameters:
			//400
			b.ReplyErrContent(w, r, http.StatusBadRequest, err.Error())
		case models.ErrRecordNotFound:
			//204 or 404?
			b.ReplyErrContent(w, r, http.StatusNotFound, err.Error())
		default:
			//500
			b.ReplyErrContent(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	//good
	render.JSON(w, r, Response{
		Status: "success",
	})
}

// GetAll list all
func (b *Building) GetAll(w http.ResponseWriter, r *http.Request) {
	data := &models.BuildingGetParams{}
	//check
	rows, err := data.GetAll(b.Storage)
	//chk
	if err != nil {
		log.Println("GETALL", err)
		//404
		b.ReplyErrContent(w, r, http.StatusNotFound, err.Error())
		return
	}
	//good
	render.JSON(w, r, Response{
		Status: "success",
		Result: rows,
		Total:  len(rows),
	})
}

// GetOne get 1 row per id
func (b *Building) GetOne(w http.ResponseWriter, r *http.Request) {
	s, _ := strconv.ParseInt(strings.TrimSpace(chi.URLParam(r, "id")), 10, 64)
	fdata := NewBuildingFormGetParams(s)
	if fdata.ID == 0 {
		//400
		b.ReplyErrContent(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	//check
	data := models.NewBuildingGetOne(fdata.ID)
	row, err := data.Get(b.Storage)
	//chk
	if err != nil {
		log.Println("GET1", err)
		//404
		b.ReplyErrContent(w, r, http.StatusNotFound, err.Error())
		return
	}
	//good
	render.JSON(w, r, Response{
		Status: "success",
		Result: row,
	})
}

// Delete remove from store
func (b *Building) Delete(w http.ResponseWriter, r *http.Request) {
	s, _ := strconv.ParseInt(strings.TrimSpace(chi.URLParam(r, "id")), 10, 64)
	data := models.NewBuildingDelete(s)
	//chk
	if data.ID == 0 {
		//400
		b.ReplyErrContent(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	//chk
	if err := data.Delete(b.Storage); err != nil {
		log.Println("DELETE", err)
		switch err {
		case models.ErrRecordNotFound:
			//404
			b.ReplyErrContent(w, r, http.StatusNotFound, err.Error())
		default:
			//500
			b.ReplyErrContent(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	//good
	render.JSON(w, r, Response{
		Status: "success",
	})
}

// ReplyErrContent send err-code/err-msg
func (b *Building) ReplyErrContent(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.Status(r, code)
	render.JSON(w, r, Response{
		Status: msg,
	})
}

// HealthCheck index page
func (b *Building) HealthCheck(w http.ResponseWriter, r *http.Request) {

	info := struct {
		Application string `json:"application"`
		BuildTime   string `json:"build"`
		Commit      string `json:"commit"`
		Release     string `json:"release"`
		Now         string `json:"now"`
	}{
		configs.Application,
		configs.BuildTime,
		configs.Commit,
		configs.Release,
		time.Now().Format(time.RFC3339),
	}

	render.JSON(w, r, info)
}
