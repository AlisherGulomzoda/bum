package populate

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	grades "bum-service/internal/service/grade-standard"
)

const (
	// FirstGrade is the first grade.
	FirstGrade = iota + 1
	// SecondGrade is the second grade.
	SecondGrade
	// ThirdGrade is the third grade.
	ThirdGrade
	// FourthGrade is the fourth grade.
	FourthGrade
	// FifthGrade is the fifth grade.
	FifthGrade
	// SixthGrade is the sixth grade.
	SixthGrade
	// SeventhGrade is the seventh grade.
	SeventhGrade
	// EighthGrade is the eighth grade.
	EighthGrade
	// NinthGrade is the eighth grade.
	NinthGrade
	// TenthGrade is the tenth grade.
	TenthGrade
	// EleventhGrade is the eleventh grade.
	EleventhGrade
	// TwelfthGrade is the twelfth grade.
	TwelfthGrade
	// ThirteenthGrade is the thirteenth grade.
	ThirteenthGrade
)

const (
	// RussianStandard is Russian standard name.
	RussianStandard = "Russian"
	// TajikStandard is Tajik standard name.
	TajikStandard = "Tajik"
	// UKStandard is UK standard name.
	UKStandard = "United Kingdom"
	// USStandard is US standard name.
	USStandard = "United States"
	// JapaneseStandard is Japanese standard name.
	JapaneseStandard = "Japanese"
)

const (
	// ElevenEduYears represents 11 education years.
	ElevenEduYears = 11
	// TwelveEduYears represents 12 education years.
	TwelveEduYears = 12
	// ThirteenEduYears represents 13 education years.
	ThirteenEduYears = 13
)

func pointerFromString(s string) *string {
	return &s
}

//nolint:funlen // it's ok
func (s *Service) addGradeStandards(ctx context.Context) error {
	gradeStandards := []grades.CreateGradeStandardArgs{
		{
			Name:           RussianStandard,
			Description:    pointerFromString("Russian Education System Grades"),
			EducationYears: ElevenEduYears,
			Grades: []grades.CreateGradeArgs{
				{
					Name:          "First Class",
					EducationYear: getPointOfInt8(FirstGrade),
				},
				{
					Name:          "Second Class",
					EducationYear: getPointOfInt8(SecondGrade),
				},
				{
					Name:          "Third Class",
					EducationYear: getPointOfInt8(ThirdGrade),
				},
				{
					Name:          "Fourth Class",
					EducationYear: getPointOfInt8(FourthGrade),
				},
				{
					Name:          "Fifth Class",
					EducationYear: getPointOfInt8(FifthGrade),
				},
				{
					Name:          "Sixth Class",
					EducationYear: getPointOfInt8(SixthGrade),
				},
				{
					Name:          "Seventh Class",
					EducationYear: getPointOfInt8(SeventhGrade),
				},
				{
					Name:          "Eighth Class",
					EducationYear: getPointOfInt8(EighthGrade),
				},
				{
					Name:          "Ninth Class",
					EducationYear: getPointOfInt8(NinthGrade),
				},
				{
					Name:          "Tenth Class",
					EducationYear: getPointOfInt8(TenthGrade),
				},
				{
					Name:          "Eleventh Class",
					EducationYear: getPointOfInt8(EleventhGrade),
				},
			},
		},
		{
			Name:           TajikStandard,
			Description:    pointerFromString("Tajik Education System Grades"),
			EducationYears: ElevenEduYears,
			Grades: []grades.CreateGradeArgs{
				{
					Name:          "First Class",
					EducationYear: getPointOfInt8(FirstGrade),
				},
				{
					Name:          "Second Class",
					EducationYear: getPointOfInt8(SecondGrade),
				},
				{
					Name:          "Third Class",
					EducationYear: getPointOfInt8(ThirdGrade),
				},
				{
					Name:          "Fourth Class",
					EducationYear: getPointOfInt8(FourthGrade),
				},
				{
					Name:          "Fifth Class",
					EducationYear: getPointOfInt8(FifthGrade),
				},
				{
					Name:          "Sixth Class",
					EducationYear: getPointOfInt8(SixthGrade),
				},
				{
					Name:          "Seventh Class",
					EducationYear: getPointOfInt8(SeventhGrade),
				},
				{
					Name:          "Eighth Class",
					EducationYear: getPointOfInt8(EighthGrade),
				},
				{
					Name:          "Ninth Class",
					EducationYear: getPointOfInt8(NinthGrade),
				},
				{
					Name:          "Tenth Class",
					EducationYear: getPointOfInt8(TenthGrade),
				},
				{
					Name:          "Eleventh Class",
					EducationYear: getPointOfInt8(EleventhGrade),
				},
			},
		},
		{
			Name:           UKStandard,
			Description:    pointerFromString("United Kingdom Education System Grades"),
			EducationYears: ThirteenEduYears,
			Grades: []grades.CreateGradeArgs{
				{
					Name:          "Year One",
					EducationYear: getPointOfInt8(FirstGrade),
				},
				{
					Name:          "Year Two",
					EducationYear: getPointOfInt8(SecondGrade),
				},
				{
					Name:          "Year Three",
					EducationYear: getPointOfInt8(ThirdGrade),
				},
				{
					Name:          "Year Four",
					EducationYear: getPointOfInt8(FourthGrade),
				},
				{
					Name:          "Year Five",
					EducationYear: getPointOfInt8(FifthGrade),
				},
				{
					Name:          "Year Six",
					EducationYear: getPointOfInt8(SixthGrade),
				},
				{
					Name:          "Year Seven",
					EducationYear: getPointOfInt8(SeventhGrade),
				},
				{
					Name:          "Year Eight",
					EducationYear: getPointOfInt8(EighthGrade),
				},
				{
					Name:          "Year Nine",
					EducationYear: getPointOfInt8(NinthGrade),
				},
				{
					Name:          "Year Ten",
					EducationYear: getPointOfInt8(TenthGrade),
				},
				{
					Name:          "Year Eleven",
					EducationYear: getPointOfInt8(EleventhGrade),
				},
				{
					Name:          "Year Twelve",
					EducationYear: getPointOfInt8(TwelfthGrade),
				},
				{
					Name:          "Year Thirteen",
					EducationYear: getPointOfInt8(ThirteenthGrade),
				},
			},
		},
		{
			Name:           USStandard,
			Description:    pointerFromString("United States Education System Grades"),
			EducationYears: TwelveEduYears,
			Grades: []grades.CreateGradeArgs{
				{
					Name:          "Kindergarten",
					EducationYear: nil,
				},
				{
					Name:          "First Grade",
					EducationYear: getPointOfInt8(FirstGrade),
				},
				{
					Name:          "Second Grade",
					EducationYear: getPointOfInt8(SecondGrade),
				},
				{
					Name:          "Third Grade",
					EducationYear: getPointOfInt8(ThirdGrade),
				},
				{
					Name:          "Fourth Grade",
					EducationYear: getPointOfInt8(FourthGrade),
				},
				{
					Name:          "Fifth Grade",
					EducationYear: getPointOfInt8(FifthGrade),
				},
				{
					Name:          "Sixth Grade",
					EducationYear: getPointOfInt8(SixthGrade),
				},
				{
					Name:          "Seventh Grade",
					EducationYear: getPointOfInt8(SeventhGrade),
				},
				{
					Name:          "Eighth Grade",
					EducationYear: getPointOfInt8(EighthGrade),
				},
				{
					Name:          "Ninth Grade",
					EducationYear: getPointOfInt8(NinthGrade),
				},
				{
					Name:          "Tenth Grade",
					EducationYear: getPointOfInt8(TenthGrade),
				},
				{
					Name:          "Eleventh Grade",
					EducationYear: getPointOfInt8(EleventhGrade),
				},
				{
					Name:          "Twelfth Grade",
					EducationYear: getPointOfInt8(TwelfthGrade),
				},
			},
		},
		{
			Name:           JapaneseStandard,
			Description:    pointerFromString("Japanese Education System Grades"),
			EducationYears: TwelfthGrade,
			Grades: []grades.CreateGradeArgs{
				{
					Name:          "First Grade (Elementary School)",
					EducationYear: getPointOfInt8(FirstGrade),
				},
				{
					Name:          "Second Grade (Elementary School)",
					EducationYear: getPointOfInt8(SecondGrade),
				},
				{
					Name:          "Third Grade (Elementary School)",
					EducationYear: getPointOfInt8(ThirdGrade),
				},
				{
					Name:          "Fourth Grade (Elementary School)",
					EducationYear: getPointOfInt8(FourthGrade),
				},
				{
					Name:          "Fifth Grade (Elementary School)",
					EducationYear: getPointOfInt8(FifthGrade),
				},
				{
					Name:          "Sixth Grade (Elementary School)",
					EducationYear: getPointOfInt8(SixthGrade),
				},
				{
					Name:          "Seventh Grade (Junior High School)",
					EducationYear: getPointOfInt8(SeventhGrade),
				},
				{
					Name:          "Eighth Grade (Junior High School)",
					EducationYear: getPointOfInt8(EighthGrade),
				},
				{
					Name:          "Ninth Grade (Junior High School)",
					EducationYear: getPointOfInt8(NinthGrade),
				},
				{
					Name:          "Tenth Grade (Senior High School)",
					EducationYear: getPointOfInt8(TenthGrade),
				},
				{
					Name:          "Eleventh Grade (Senior High School)",
					EducationYear: getPointOfInt8(EleventhGrade),
				},
				{
					Name:          "Twelfth Grade (Senior High School)",
					EducationYear: getPointOfInt8(TwelfthGrade),
				},
			},
		},
	}

	for _, gradeStandard := range gradeStandards {
		g, err := s.gradesService().CreateGradeStandard(ctx, gradeStandard)
		if err != nil {
			return fmt.Errorf("failed to create grade standard: %w", err)
		}

		s.gradesList[g.Name] = g
		s.gradesList[g.ID.String()] = g
	}

	return nil
}

func (s *Service) getGradeByName(gradeStandard string, eduYear int8) domain.Grade {
	g := s.gradesList[gradeStandard]
	for _, grade := range g.Grades {
		if grade.EducationYear != nil && *grade.EducationYear == eduYear {
			return grade
		}
	}

	return domain.Grade{}
}

func (s *Service) getGradeByID(gradeStandardID uuid.UUID, eduYear int8) domain.Grade {
	g := s.gradesList[gradeStandardID.String()]
	for _, grade := range g.Grades {
		if grade.EducationYear != nil && *grade.EducationYear == eduYear {
			return grade
		}
	}

	return domain.Grade{}
}

func (s *Service) getGradeStandardByID(gradeStandardID uuid.UUID) domain.GradeStandard {
	return s.gradesList[gradeStandardID.String()]
}
