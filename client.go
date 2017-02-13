package goxbdm

import (
	"errors"
	"strconv"
	"strings"

	"fmt"

	"regexp"

	"github.com/0xdeafcafe/go-xbdm/clients"
)

// Client ..
type Client struct {
	tcpClient *clients.TCPClient

	consoleName string
	ConsoleIP   string
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
		return "", err
	}

	resp, err := client.tcpClient.ReadString()
	if err != nil {
		return "", err
	}

	_, successful, message := validateResponse(resp)
	if !successful {
		return "", errors.New(message)
	}

	return message, nil
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

// ConsoleName returns the name of the console at the time the connection was established.
func (client *Client) ConsoleName() string {
	return client.consoleName
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

	// Create client
	client := &Client{
		tcpClient: tcpClient,
		ConsoleIP: xboxIP,
	}

	// Store debug console name
	dbgName, err := client.DebugName()
	if err != nil {
		return nil, err
	}

	client.consoleName = dbgName
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

// validateResponse validates a string response from the Xbox.
func validateResponse(str string) (status int, successful bool, message string) {
	r := regexp.MustCompile(`(?P<status>[\d]{3})-\s(?P<message>.*)`)
	matches := r.FindStringSubmatch(str)
	names := r.SubexpNames()

	if len(matches) != 3 {
		return -1, false, str
	}

	status = -1
	message = ""
	for i := range matches {
		switch names[i] {
		case "status":
			status, _ = strconv.Atoi(matches[i])
			break
		case "message":
			message = matches[i]
			break
		}
	}

	return status, status >= 200 && status <= 299, message
}
