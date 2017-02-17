package models

import "fmt"
import "strings"

// Drive defines the structure of the `drivelist` response
type Drive struct {
	Name string
}

// NewDrive ..
func NewDrive(body string) *Drive {
	format := `drivename=%s`
	var name string
	fmt.Sscanf(body, format, &name)

	return &Drive{
		Name: strings.Trim(name, `"`),
	}
}
