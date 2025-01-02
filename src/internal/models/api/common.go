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

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) GetLimit() int {
	return p.PageSize
}

func ConvertToPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	// Handle default page size
	pageSize = max(pageSize, DefaultPageSize)

	// Calculate total pages
	totalPages := total / int64(pageSize)

	// Add an extra page if there are remaining items
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}
}
