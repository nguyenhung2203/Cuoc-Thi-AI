package pagination

import (
	"net/http"
	"strconv"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Params holds the parsed pagination and sort query parameters.
type Params struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
}

// Meta holds the pagination metadata returned in API responses.
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// FromRequest parses page, page_size, sort_by, and sort_dir from query params.
// Defaults: page=1, page_size=20 (max 100).
func FromRequest(r *http.Request) Params {
	q := r.URL.Query()

	page := parseInt(q.Get("page"), 1)
	if page < 1 {
		page = 1
	}

	pageSize := parseInt(q.Get("page_size"), defaultPageSize)
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	sortDir := q.Get("sort_dir")
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
		SortBy:   q.Get("sort_by"),
		SortDir:  sortDir,
	}
}

// Offset returns the SQL OFFSET value for the current page.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// CalcTotalPages computes the number of pages needed for total items.
func CalcTotalPages(total, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	return pages
}

// parseInt parses s as an int; returns fallback on empty string or error.
func parseInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}
