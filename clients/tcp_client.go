package clients

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

// TCPClient defines the structure of the TCP client.
type TCPClient struct {
	Host       string
	connection net.Conn
}

var connectionPulseMessage = []byte{0x32, 0x30, 0x31, 0x2d, 0x20, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x0d, 0x0a}

// Write sends a byte slice to the xbox.
func (client *TCPClient) Write(data []byte) {
	client.connection.Write(data)
}

// WriteString sends a string to the xbox. If `addLineEndings` is true, it will append
// `\r\n` to the message.
func (client *TCPClient) WriteString(message string, addLineEndings bool) {
	if addLineEndings {
		client.connection.Write([]byte(fmt.Sprintf("%s\r\n", message)))
	} else {
		client.connection.Write([]byte(message))
	}
}

// Read reads the pending response from the xbox.
func (client *TCPClient) Read(bufferSize int) []byte {
	eom := []byte{0x0d, 0x0a}                         // \r\n
	multilinePrefix := []byte{0x32, 0x30, 0x32, 0x2d} /* 202- */
	message := make([]byte, 0)

	i := 0
	for {
		// Read buffer in
		buffer := make([]byte, bufferSize)
		client.connection.Read(buffer)
		message = append(message, buffer...)

		// If first read, and response is multi-line, update end-of-message slice
		if i == 0 && bytes.HasPrefix(buffer, multilinePrefix) {
			eom = []byte{0x2e, 0x0d, 0x0a} // .\r\n
		}

		// Check if message has ended
		trimmedBuf := bytes.TrimRight(message, "\x00")
		if bytes.HasSuffix(trimmedBuf, eom) {
			return trimmedBuf
		}

		i++
	}
}

// ReadString reads the pending response from the xbox into a string.
func (client *TCPClient) ReadString() string {
	return string(client.Read())
}

// Close ends the open TCP connection.
func (client *TCPClient) Close() {
	// Wish our loved ones goodbye
	client.WriteString("bye", true)
	client.connection.Close()
}

// NewTCPClientWithPort creates a new TCPClient with an Xbox's IP and listening port.
func NewTCPClientWithPort(xboxIP string, port int) (*TCPClient, error) {
	client := &TCPClient{
		Host: fmt.Sprintf("%s:%d", xboxIP, port),
	}

	// Connect to Xbox
	conn, err := net.Dial("tcp", client.Host)
	if err != nil {
		return nil, err
	}

	// Set connection
	client.connection = conn

	// Check connection pulse exists
	if !bytes.Equal(client.Read(), connectionPulseMessage) {
		return nil, errors.New("connection pulse message failed")
	}

	return client, nil
}

// responseIsError checks to see if the response message is valid of an error
func responseIsError(str string) bool {
	r := regexp.MustCompile(`(?P<code>[\d]{3})`)
	matches := r.FindStringSubmatch(str)
	names := r.SubexpNames()
	if len(matches) < 2 {
		return false
	}

	code := -1

	for i := range matches {
		if names[i] == "code" {
			code, _ = strconv.Atoi(matches[i])
		}
	}

	return code <= 299 && code >= 200
}
