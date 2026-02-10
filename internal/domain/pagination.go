package domain

// Pagination is structure of pagination.
type Pagination struct {
	Limit  int
	Offset int
}

const (
	// PaginationDefaultLimit is the default limit for pagination.
	PaginationDefaultLimit = 10
	// PaginationDefaultPage is the default page for pagination.
	PaginationDefaultPage = 1
)

// NewPagination creates a new Pagination domain.
func NewPagination(page, perPage int) Pagination {
	return Pagination{
		Limit:  perPage,
		Offset: (page - 1) * perPage,
	}
}
