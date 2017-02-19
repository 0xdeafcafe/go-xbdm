package goxbdm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/0xdeafcafe/go-xbdm/models"
)

// rebootType is a custom type for storing different types of reboot.
type rebootType uint8

// freezeState is a custom type for storing the different states of console freezing.
type freezeState uint8

const (
	screenshotHeaderFormat         = "pitch=0x%x width=0x%x height=0x%x format=0x%x offsetx=0x00000000 offsety=0x00000000, framebuffersize=0x%x"
	rebootTitleToActiveTitleFormat = `magicboot title=%s directory=%s`

	// RebootTitle defines the enum for rebooting to the Developer Dashboard.
	RebootTitle rebootType = iota

	// RebootTitleToActiveTitle defines the enum for rebooting to the currently active
	// title.
	RebootTitleToActiveTitle

	// RebootCold defines the enum for turning the kit off and then back on.
	RebootCold

	// FrozenState is when all the threads are resumed as normal.
	FrozenState freezeState = iota

	// UnfrozenState is when all the threads are paused.
	UnfrozenState
)

// Reboot the Xbox console.
func (client *Client) Reboot(rebootType rebootType) error {
	switch {
	case rebootType == RebootTitle:
		_, err := client.SendCommand("magicboot")
		return err

	case rebootType == RebootTitleToActiveTitle:
		info, err := client.RunningXBEInfo()
		if err != nil {
			return err
		}

		// Get directory of the currently running title
		name := info.Name
		titleDirectory := fmt.Sprintf(`%s"`, name[0:strings.LastIndex(name, "\\")])

		// Tell Xbox which title to launch
		_, err = client.SendCommand(fmt.Sprintf(rebootTitleToActiveTitleFormat, name, titleDirectory))
		return err

	case rebootType == RebootCold:
		_, err := client.SendCommand("magicboot COLD")
		return err

	default:
		panic("invalid reboot type")
	}
}

// RunningXBEInfo gets the xbeinfo of the currently running title.
func (client *Client) RunningXBEInfo() (*models.XBEInfo, error) {
	_, err := client.SendCommand("xbeinfo running")
	if err != nil {
		return nil, err
	}

	body, err := client.ReadMultilineResponse()
	if err != nil {
		return nil, err
	}

	return models.NewXBEInfo(body), nil
}

// Screenshot dumps the frame buffer of the Xbox.
func (client *Client) Screenshot() ([]byte, error) {
	resp, err := client.SendCommand("screenshot")
	if err != nil {
		return nil, err
	}
	if resp != "203- binary response follows" {
		return nil, errors.New(resp)
	}

	// Read header values
	header, _ := client.tcpClient.ReadString()
	var pitch, width, height, format, frameBufferSize int
	fmt.Sscanf(header, screenshotHeaderFormat, &pitch, &width, &height, &format, &frameBufferSize)
	buf := bytes.NewBuffer(nil)
	buf.Grow(frameBufferSize)

	_, err = io.CopyN(buf, client.tcpClient.Reader, int64(frameBufferSize))
	if err != nil {
		log.Fatal("Reading image data failed: ", err)
	}

	// Deswizzle this thing
	data := buf.Bytes()
	for i := 0; i < pitch*height; i += pitch / width {
		data[i], data[i+2] = data[i+2], data[i]
	}

	return data, nil
}

// ChangeFreezeState set's the freeze state of the Xbox.
func (client *Client) ChangeFreezeState(freezeState freezeState) (err error) {
	switch freezeState {
	case FrozenState:
		_, err = client.SendCommand("freeze")
		break
	case UnfrozenState:
		_, err = client.SendCommand("unfreeze")
		break
	}
	return
}
