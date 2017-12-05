package paging

import (
	"database/sql"
	"fmt"
	"strings"
)

var (
	DB *sql.DB

	DefaultPageQuery = &PageQuery{
		NumCurrentPage: 0,
		NumPerPage:     20,
		OrderBy:        nil,
	}
)

type PageQuery struct {
	NumCurrentPage int
	NumPerPage     int
	OrderBy        []string
}

func (pq *PageQuery) ToString() string {
	s := ""
	if len(pq.OrderBy) != 0 {
		s += " ORDER BY " + strings.Join(pq.OrderBy, ",")
	}
	s += fmt.Sprintf(" LIMIT %d OFFSET %d", pq.NumPerPage, pq.NumCurrentPage*pq.NumPerPage)
	return s
}
