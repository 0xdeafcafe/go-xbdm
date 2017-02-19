package goxbdm

import "fmt"

func (client *Client) SetMemory(address int64, data string) (string, error) {
	return client.SendCommand(fmt.Sprintf(`setmem addr=0x%x data=%s`, address, data))
}

func (client *Client) GetMemory(address, length int64) ([]byte, error) {
	_, err := client.SendCommand(fmt.Sprintf(`getmemex addr=0x%x length=%x`, address, length))
	if err != nil {
		return nil, err
	}

	data, _, err := client.tcpClient.Read(length)
	return data, err
}
