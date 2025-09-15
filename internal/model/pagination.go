package model

// PaginationMetadata contém informações sobre a paginação
type PaginationMetadata struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Data     interface{}        `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

// CalculatePagination calcula os metadados de paginação
func CalculatePagination(page, limit int, total int64) PaginationMetadata {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Limite máximo
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMetadata{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
