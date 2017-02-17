package goxbdm

import "fmt"

func (client *Client) SetMemory(address int64, data string) (string, error) {
	return client.SendCommand(fmt.Sprintf(`setmem addr=0x%x data=%s`, address, data))
}
