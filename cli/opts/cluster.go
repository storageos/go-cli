package opts

import "fmt"

// ValidateClusterSize validates whether the input string is a valid cluster size.
func ValidateClusterSize(val int) (int, error) {
	switch val {
	case 1, 3, 5, 7:
		return val, nil
	default:
		return 0, fmt.Errorf("cluster consensus size must be one of %d, %d, %d, or %d (minimum 3 for production)", 1, 3, 5, 7)
	}
}
