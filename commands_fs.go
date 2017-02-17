package goxbdm

import "github.com/0xdeafcafe/go-xbdm/models"
import "fmt"

// ListDrives ..
func (client *Client) ListDrives() ([]*models.Drive, error) {
	_, err := client.SendCommand("drivelist")
	if err != nil {
		return nil, err
	}

	body, err := client.ReadMultilineResponse()
	if err != nil {
		return nil, err
	}

	drives := make([]*models.Drive, len(body))
	for i := 0; i < len(body); i++ {
		drives[i] = models.NewDrive(body[i])
	}

	return drives, nil
}

// ListDirectory ..
func (client *Client) ListDirectory(dir string) ([]*models.DirectoryItem, error) {
	_, err := client.SendCommand(fmt.Sprintf(`dirlist name="%s"`, dir))
	if err != nil {
		return nil, err
	}

	body, err := client.ReadMultilineResponse()
	if err != nil {
		return nil, err
	}

	directoryItems := make([]*models.DirectoryItem, len(body))
	for i := 0; i < len(body); i++ {
		directoryItems[i] = models.NewDirectoryItem(body[i])
	}

	return directoryItems, nil
}
