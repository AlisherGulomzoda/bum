package request

import "time"

// Pagination is structure of pagination request.
type Pagination struct {
	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=10" binding:"min=1,max=500"`
}

// ListFilter is structure of list request.
type ListFilter struct {
	SortOrder string `form:"sort_order,default=DESC" binding:"oneof=DESC ASC"`
	Pagination
}

// DateFilter date time filter for embedding.
type DateFilter struct {
	From *string `form:"date_from" binding:"omitempty" time_format:"2006-01-02"`
	Till *string `form:"date_till" binding:"omitempty" time_format:"2006-01-02"`
}

// DateFrom converts date from to time.Time.
//
//nolint:errcheck // err always is nil because we already checked it before.
func (f DateFilter) DateFrom() *time.Time {
	if f.From == nil {
		return nil
	}

	t, _ := time.ParseInLocation(time.DateOnly, *f.From, time.UTC)

	return &t
}

// DateTill converts date till to time.Time.
//
//nolint:errcheck // err always is nil because we already checked it before.
func (f DateFilter) DateTill() *time.Time {
	if f.Till == nil {
		return nil
	}

	t, _ := time.ParseInLocation(time.DateOnly, *f.Till, time.UTC)

	return &t
}
