package tools

import (
	"net/http"
	"strconv"
	"strings"
)

// PagingParams paging details from req-parameters
type PagingParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total,omitempty"`
}

// NewPagingParams parse paging
func NewPagingParams(r *http.Request) *PagingParams {

	paging := &PagingParams{
		Page:  1,
		Limit: 10,
	}

	if p, err := strconv.Atoi(strings.TrimSpace(
		r.URL.Query().Get("page"))); err == nil && p > 0 {
		paging.Page = p
	}

	if p, err := strconv.Atoi(strings.TrimSpace(
		r.URL.Query().Get("limit"))); err == nil && p > 0 {
		paging.Limit = p
	}
	return paging
}
