package model

import (
	"fmt"
	"math"
	"strings"
)

type Pageable struct {
	Page      int
	Size      int
	Sort      []string
	Direction string
}

// NewPageable 建立一個預設值的安全分頁物件
func NewPageable(page, size *int, direction string, sort ...string) Pageable {
	if *page < 1 {
		*page = 1
	}
	if *size < 1 {
		*size = 10
	}
	if *size > 500 { // 防止一次拉太多資料
		*size = 500
	}

	p := Pageable{
		Page: *page,
		Size: *size,
		Direction: direction,
	}

	if len(sort) > 0 {
		p.Sort = sort
	}

	return p
}

func (p Pageable) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p Pageable) Limit() int {
	return p.Size
}

func (p Pageable) TotalPages(total int64) int64 {
	if total == 0 {
		return 0
	}
	return int64(math.Ceil(float64(total) / float64(p.Size)))
}

func (p Pageable) IsASC() bool {
	return p.Direction == "asc" || p.Direction == "ASC"
}

func (p Pageable) OrderBySQL() string {
	if len(p.Sort) == 0 {
		return "updated_at DESC"
	}

	dir := "DESC"
	if p.IsASC() {
		dir = "ASC"
	}

	var orders []string
	for _, s := range p.Sort {
		if isValidColumn(s) {
			orders = append(orders, fmt.Sprintf("%s %s", s, dir))
		}
	}

	if len(orders) == 0 {
		return "updated_at DESC"
	}

	return strings.Join(orders, ", ")
}

func isValidColumn(col string) bool {
	if col == "" {
		return false
	}
	for _, ch := range col {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}
