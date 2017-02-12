package clients

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

// TCPClient defines the structure of the TCP client.
type TCPClient struct {
	Host       string
	connection net.Conn

	Reader *bufio.Reader
	Writer *bufio.Writer
}

const (
	messageSuffix = "\r\n"
)

var connectionPulseMessage = `201- connected`

// Write sends a byte slice to the xbox.
func (client *TCPClient) Write(data []byte) (int, error) {
	n, err := client.Writer.Write(data)
	if err != nil {
		return n, err
	}

	err = client.Writer.Flush()
	if err != nil {
		return -1, err
	}

	return n, nil
}

// WriteString sends a command to the Xbox. This function will suffix the command with
// `\r\n`.
func (client *TCPClient) WriteString(message string) (int, error) {
	return client.Write([]byte(message + messageSuffix))
}

// ReadString reads a string response from the Xbox.
func (client *TCPClient) ReadString() (string, error) {
	str, err := client.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(str, "\r\n"), nil
}

// Read reads the pending response from the xbox.
func (client *TCPClient) Read(len int) ([]byte, int, error) {
	buffer := make([]byte, len)
	n, err := client.Reader.Read(buffer)
	return buffer, n, err
}

// Close ends the open TCP connection.
func (client *TCPClient) Close() {
	// Wish our loved ones goodbye
	client.WriteString("bye")
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
	client.Reader = bufio.NewReader(client.connection)
	client.Writer = bufio.NewWriter(client.connection)

	// Check connection pulse exists
	pulse, err := client.ReadString()
	if pulse != connectionPulseMessage {
		return nil, errors.New("connection pulse message failed")
	}

	return client, nil
}
