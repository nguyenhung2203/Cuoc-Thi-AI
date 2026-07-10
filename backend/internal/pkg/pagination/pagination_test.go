package pagination

import (
	"net/http"
	"net/url"
	"testing"
)

func TestFromRequest_Defaults(t *testing.T) {
	r := &http.Request{URL: &url.URL{RawQuery: ""}}
	p := FromRequest(r)
	if p.Page != 1 {
		t.Fatalf("expected page 1, got %d", p.Page)
	}
	if p.PageSize != 20 {
		t.Fatalf("expected page_size 20, got %d", p.PageSize)
	}
	if p.SortDir != "asc" {
		t.Fatalf("expected asc, got %s", p.SortDir)
	}
}

func TestFromRequest_Custom(t *testing.T) {
	r := &http.Request{URL: &url.URL{RawQuery: "page=3&page_size=50&sort_by=created_at&sort_dir=desc"}}
	p := FromRequest(r)
	if p.Page != 3 {
		t.Fatalf("expected page 3, got %d", p.Page)
	}
	if p.PageSize != 50 {
		t.Fatalf("expected page_size 50, got %d", p.PageSize)
	}
	if p.SortBy != "created_at" {
		t.Fatalf("expected created_at, got %s", p.SortBy)
	}
	if p.SortDir != "desc" {
		t.Fatalf("expected desc, got %s", p.SortDir)
	}
}

func TestFromRequest_MaxPageSize(t *testing.T) {
	r := &http.Request{URL: &url.URL{RawQuery: "page_size=999"}}
	p := FromRequest(r)
	if p.PageSize > 100 {
		t.Fatalf("expected max 100, got %d", p.PageSize)
	}
}

func TestFromRequest_InvalidPage(t *testing.T) {
	r := &http.Request{URL: &url.URL{RawQuery: "page=-1&page_size=0"}}
	p := FromRequest(r)
	if p.Page != 1 {
		t.Fatalf("expected page 1, got %d", p.Page)
	}
	if p.PageSize != 20 {
		t.Fatalf("expected page_size 20, got %d", p.PageSize)
	}
}

func TestOffset(t *testing.T) {
	p := Params{Page: 3, PageSize: 10}
	off := p.Offset()
	if off != 20 {
		t.Fatalf("expected 20, got %d", off)
	}
}

func TestCalcTotalPages(t *testing.T) {
	tests := []struct {
		total, size, want int
	}{
		{0, 20, 0},
		{100, 20, 5},
		{101, 20, 6},
		{1, 20, 1},
		{0, 0, 0},
	}
	for _, tt := range tests {
		got := CalcTotalPages(tt.total, tt.size)
		if got != tt.want {
			t.Errorf("CalcTotalPages(%d, %d) = %d; want %d", tt.total, tt.size, got, tt.want)
		}
	}
}

func TestParseInt(t *testing.T) {
	if got := parseInt("42", 0); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if got := parseInt("", 10); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
	if got := parseInt("abc", 5); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}
