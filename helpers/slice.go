package helpers

// StrSliceContains checks if an array contains a string
func StrSliceContains(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}

	return false
}
