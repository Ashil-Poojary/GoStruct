package utils

import "strconv"

// StringToInt64
// Converts a string to int64, returns 0 and error if conversion fails
func StringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil // Return 0 for empty string
	}
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err // Return error if conversion fails
	}
	return value, nil
}
