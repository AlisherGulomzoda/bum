package request

import "github.com/gin-gonic/gin"

const (
	eduOrganizationIDPathVar = "edu_organization_id" // eduOrganizationIDParam is educational organization id param.
	schoolIDPathVar          = "school_id"           // schoolIDParam is school id param.
	schoolSubjectIDPathVar   = "school_subject_id"   // schoolIDParam is school subject id param.
	auditoriumIDPathVar      = "auditorium_id"       // auditoriumIDPathVar is school auditorium id param.
	directorIDPathVar        = "director_id"         // directorIDParam is director id param.
	ownerIDPathVar           = "owner_id"            // ownerIDPathVar is owner id param.
	guardianIDPathVar        = "guardian_id"         // guardianIDParam is guardian id param.
	userIDPathVar            = "user_id"             // userIDParam is user id param.
	headmasterIDPathVar      = "headmaster_id"       // headmasterIDParam is headmaster id param.
	subjectIDPathVar         = "subject_id"          // subjectIDParam is subject id param.
	teacherIDPathVar         = "teacher_id"          // teacherIDParam is teacher id param.
	gradeStandardIDPathVar   = "grade_standard_id"   // gradeStandardIDPathVar is grade standard id param.
	groupIDPathVar           = "group_id"            // groupIDPathVar is group id param.
	groupSubjectIDPathVar    = "group_subject_id"    // groupSubjectIDPathVar is group subject id param.
	studyPlanIDPathVar       = "study_plan_id"       // studyPlanIDPathVar is study plan id param.
	studyPlanStatusPathVar   = "study_plan_status"   // studyPlanStatusPathVar is study plan status param.
	studentIDPathVar         = "student_id"          // studentIDPathVar is student id param
	markIDPathVar            = "mark_id"             // markIDPathVar is mark id param
)

// GetEduOrganizationPathVar gets edu organization id from path variable.
func GetEduOrganizationPathVar(c *gin.Context) string {
	return c.Param(eduOrganizationIDPathVar)
}

// GetSchoolIDPathVar gets edu school id from path variable.
func GetSchoolIDPathVar(c *gin.Context) string {
	return c.Param(schoolIDPathVar)
}

// GetSchoolSubjectIDPathVar gets school subject id from path variable.
func GetSchoolSubjectIDPathVar(c *gin.Context) string {
	return c.Param(schoolSubjectIDPathVar)
}

// GetSchoolAuditoriumIDPathVar gets school auditorium id from path variable.
func GetSchoolAuditoriumIDPathVar(c *gin.Context) string {
	return c.Param(auditoriumIDPathVar)
}

// GetDirectorIDPathVar gets director id from path variable.
func GetDirectorIDPathVar(c *gin.Context) string {
	return c.Param(directorIDPathVar)
}

// GetOwnerIDPathVar gets owner id from path variable.
func GetOwnerIDPathVar(c *gin.Context) string {
	return c.Param(ownerIDPathVar)
}

// GetGuardianIDPathVar gets guardian id from path variable.
func GetGuardianIDPathVar(c *gin.Context) string {
	return c.Param(guardianIDPathVar)
}

// GetUserIDPathVar gets user id from path variable.
func GetUserIDPathVar(c *gin.Context) string {
	return c.Param(userIDPathVar)
}

// GetHeadmasterIDPathVar gets headmaster id from path variable.
func GetHeadmasterIDPathVar(c *gin.Context) string {
	return c.Param(headmasterIDPathVar)
}

// GetSubjectIDPathVar gets subject_id from path variable.
func GetSubjectIDPathVar(c *gin.Context) string {
	return c.Param(subjectIDPathVar)
}

// GetTeacherIDPathVar gets teacher id from path variable.
func GetTeacherIDPathVar(c *gin.Context) string {
	return c.Param(teacherIDPathVar)
}

// GetGradeStandardIDPathVar gets grade standard id from path variable.
func GetGradeStandardIDPathVar(c *gin.Context) string {
	return c.Param(gradeStandardIDPathVar)
}

// GetGroupIDPathVar gets group id from path variable.
func GetGroupIDPathVar(c *gin.Context) string {
	return c.Param(groupIDPathVar)
}

// GetGroupSubjectIDPathVar gets group subject id from path variable.
func GetGroupSubjectIDPathVar(c *gin.Context) string {
	return c.Param(groupSubjectIDPathVar)
}

// GetStudyPlanIDPathVar gets study plan id from path variable.
func GetStudyPlanIDPathVar(c *gin.Context) string {
	return c.Param(studyPlanIDPathVar)
}

// GetStudyPlanStatusPathVar gets study plan status from path variable.
func GetStudyPlanStatusPathVar(c *gin.Context) string {
	return c.Param(studyPlanStatusPathVar)
}

// GetStudentIDPathVar gets student if from path variable.
func GetStudentIDPathVar(c *gin.Context) string { return c.Param(studentIDPathVar) }

// GetMarkIDPathVar gets mark id from path variable.
func GetMarkIDPathVar(c *gin.Context) string { return c.Param(markIDPathVar) }
