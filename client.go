package goxbdm

import (
	"strconv"
	"strings"

	"fmt"

	"github.com/0xdeafcafe/go-xbdm/clients"
)

// Client ..
type Client struct {
	tcpClient *clients.TCPClient
}

const (
	defaultXboxTCPPort = 730
)

// DebugName returns the debug name of the Xbox.
func (client *Client) DebugName() (string, error) {
	return client.SendCommand("dbgname")
}

// SendCommand sends a text command to the Xbox.
func (client *Client) SendCommand(command string) (string, error) {
	_, err := client.tcpClient.WriteString(command)
	if err != nil {
		panic(err)
	}

	return client.tcpClient.ReadString()
}

// ReadMultilineResponse reads the body of a multiline response and returns it.
func (client *Client) ReadMultilineResponse() (string, error) {
	lines := make([]string, 0)
	for {
		str, err := client.tcpClient.ReadString()
		if err != nil {
			return strings.Join(lines, " "), err
		}

		if str == "." {
			return strings.Join(lines, " "), nil
		}

		lines = append(lines, str)
	}
}

// ParseSpaceSeparatedValues parses a space-separated list of key-value mappings.
func (client *Client) ParseSpaceSeparatedValues(str string) map[string]string {
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

// Close ends the connection with the Xbox.
func (client *Client) Close() {
	client.tcpClient.Close()
}

// NewXBDMClient creates a new XBDM client.
func NewXBDMClient(xboxIP string) (*Client, error) {
	return NewXBDMClientWithPort(xboxIP, defaultXboxTCPPort)
}

// NewXBDMClientWithPort creates a new XBDM client with a custom port.
func NewXBDMClientWithPort(xboxIP string, port int) (*Client, error) {
	tcpClient, err := clients.NewTCPClientWithPort(xboxIP, port)
	if err != nil {
		return nil, err
	}

	client := &Client{
		tcpClient: tcpClient,
	}

	return client, nil
}

// parseMultilineResponse reads the space separated values into a map.
func parseMultilineResponse(str string) (map[string]string, error) {
	body := make(map[string]string)

	// Remove header and eom
	splitStr := strings.Split(str, "\r\n")
	dataLines := splitStr[1 : len(splitStr)-2]

	// Join array with spaces, then add an extra space at the end to save last entry
	dataStr := strings.Join(dataLines, " ") + " "
	currentKey := ""
	currentValue := ""
	isInsideQuote := false
	onKey := true
	for _, char := range dataStr {
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

	return body, nil
}

// convertHexToInt converts a hex string with the prefix `0x` into an int64.
func convertHexToInt(hexStr string) int64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int64(result)
}
