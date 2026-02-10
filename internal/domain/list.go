package domain

// ListFilter is filter for list.
type ListFilter struct {
	SortOrder string
	Pagination
}

const (
	// SortOrderDesc is desc sort order.
	SortOrderDesc = "DESC"
	// SortOrderASC is asc sort order.
	SortOrderASC = "ASC"
)

// NewListFilter creates a new List Filter domain.
func NewListFilter(sortOrder string, pagination Pagination) ListFilter {
	return ListFilter{
		SortOrder:  sortOrder,
		Pagination: pagination,
	}
}
