package goxbdm

import (
	"fmt"
	"strings"
)

// RebootType ..
type RebootType uint8

const (
	rebootTitleToActiveTitleFormat = `magicboot title=%s directory=%s`

	// RebootTitle defines the enum for rebooting to the Developer Dashboard.
	RebootTitle RebootType = iota

	// RebootTitleToActiveTitle defines the enum for rebooting to the currently active
	// title.
	RebootTitleToActiveTitle

	// RebootCold defines the enum for turning the kit off and then back on.
	RebootCold
)

// Reboot the Xbox console.
func (client *Client) Reboot(rebootType RebootType) error {
	switch {
	case rebootType == RebootTitle:
		client.tcpClient.WriteString("magicboot ", true)
		_, err := client.tcpClient.ReadString()
		return err

	case rebootType == RebootTitleToActiveTitle:
		client.tcpClient.WriteString("xbeinfo running ", true)
		info, err := client.tcpClient.ReadString()
		if err != nil {
			return err
		}

		// Parse response, and read out body
		body, _ := parseMultilineResponse(info)
		title := body["name"]

		// Split by last `\\` in title to get title directory
		titleDirectory := fmt.Sprintf(`%s"`, title[0:strings.LastIndex(title, "\\")])

		// Tell xbox what's gucci
		client.tcpClient.WriteString(fmt.Sprintf(rebootTitleToActiveTitleFormat, title, titleDirectory), true)
		_, err = client.tcpClient.ReadString()
		return err

	case rebootType == RebootCold:
		client.tcpClient.WriteString("magicboot  COLD", true)
		_, err := client.tcpClient.ReadString()
		return err

	default:
		panic("invalid reboot type")
	}
}
