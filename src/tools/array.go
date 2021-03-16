package tools

// InArrayStrings checker
func InArrayStrings(data []string, input string) bool {
	for _, i := range data {
		if i == input {
			return true
		}
	}

	return false
}
