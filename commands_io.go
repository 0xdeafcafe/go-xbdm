package goxbdm

import "strings"
import "fmt"

// ListDrives ..
func (client *Client) ListDrives() ([]string, error) {
	_, err := client.SendCommand("drivelist\r\n")
	if err != nil {
		return nil, err
	}

	body, err := client.ReadMultilineResponse()
	if err != nil {
		return nil, err
	}

	drives := make([]string, 0)
	for _, str := range strings.Split(body, " ") {
		format := `drivename="%s"`
		var driveName string
		fmt.Sscanf(str, format, &driveName)
		drives = append(drives, driveName)
	}

	return drives, nil
}
