package helpers

import "fmt"

// ParseSpaceSeparatedValues parses a space-separated list of key-value mappings.
func ParseSpaceSeparatedValues(str string) map[string]string {
	body := make(map[string]string)

	// Add space to the end of the string, to trigger the saveing logic
	str += " "

	currentKey := ""
	currentValue := ""
	isInsideQuote := false
	onKey := true
	for _, char := range str {
		// If we detect a space, and we aren't inside a double quote, switch to `onKey`
		if char == ' ' && !isInsideQuote {
			// Before we switch we need to save the read values to the `body` map and
			// reset the `currentKey` and `currentValue` variables
			body[currentKey] = currentValue
			currentKey = ""
			currentValue = ""
			onKey = true
			continue
		}

		// If we find an equals, and we aren't inside a double quote, switch to `!onKey`
		if char == '=' && !isInsideQuote {
			onKey = false
			continue
		}

		// If we find a double quote, switch between inside and outside
		if char == '"' {
			isInsideQuote = !isInsideQuote
		}

		// Save the value to the relevant part
		if onKey {
			currentKey = fmt.Sprintf("%s%s", currentKey, string(char))
		} else {
			currentValue = fmt.Sprintf("%s%s", currentValue, string(char))
		}
	}

	return body
}
