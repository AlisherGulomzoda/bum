package request

import "github.com/gin-gonic/gin"

const (
	schoolIDHeaderVar = "school_id" // schoolIDParam is school id param.
)

// GetSchoolIDHeader gets edu school id from header variable.
func GetSchoolIDHeader(c *gin.Context) string {
	return c.GetHeader(schoolIDHeaderVar)
}
