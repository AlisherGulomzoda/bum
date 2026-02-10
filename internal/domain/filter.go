package domain

import (
	"time"
)

// DateFilter date time filter for embedding.
type DateFilter struct {
	DateFrom *time.Time
	DateTill *time.Time
}

// NewDateFilter creates a new DateFilter domain.
func NewDateFilter(dateFrom, dateTill *time.Time) DateFilter {
	var utcDateFrom, utcDateTill *time.Time

	if dateFrom != nil {
		t := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
		utcDateFrom = &t
	}

	if dateTill != nil {
		t := time.Date(dateTill.Year(), dateTill.Month(), dateTill.Day(), 0, 0, 0, 0, time.UTC)
		utcDateTill = &t
	}

	return DateFilter{DateFrom: utcDateFrom, DateTill: utcDateTill}
}
