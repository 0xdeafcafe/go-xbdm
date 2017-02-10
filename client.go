package goxbdm

import (
	"github.com/0xdeafcafe/go-xbdm/clients"
)

// Client ..
type Client struct {
	tcpClient *clients.TCPClient
}

const (
	defaultXboxTCPPort = 730
)

// DebugName returns the debug name of the Xbox
func (client *Client) DebugName() (string, error) {
	client.tcpClient.WriteString("dbgname", true)
	return client.tcpClient.ReadString()
}

// NewXBDMClient ..
func NewXBDMClient(xboxIP string) (*Client, error) {
	return NewXBDMClientWithPort(xboxIP, defaultXboxTCPPort)
}

// NewXBDMClientWithPort ..
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
