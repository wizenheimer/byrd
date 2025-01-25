package models

const (
	DefaultPageSize = 10
)

// PaginationParams is the parameters for paginating a list of items
type PaginationParams struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	SortDir  string `query:"sort_dir"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) GetLimit() int {
	return p.PageSize
}
