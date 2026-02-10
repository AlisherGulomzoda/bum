package utils

import (
	"fmt"
	"strings"
)

// GetStrValueFromArgument gets the string value from the passed arguments.
func GetStrValueFromArgument[T any](input *T) *string {
	if input != nil {
		str := fmt.Sprint(*input)
		if trimmedStr := strings.TrimSpace(str); trimmedStr != "" {
			return &trimmedStr
		}
	}

	return nil
}
