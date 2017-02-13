package helpers

import "fmt"

// ConvertHexToInt64 converts a 0x preffixed hex string to an int64.
func ConvertHexToInt64(hexStr string) int64 {
	var val int64
	fmt.Sscanf("0x%x", hexStr, &val)
	return val
}
