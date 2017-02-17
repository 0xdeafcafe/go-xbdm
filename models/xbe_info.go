package models

import (
	"strings"

	"github.com/0xdeafcafe/go-xbdm/helpers"
)

// XBEInfo defines the structure of the `xbeinfo running` response
type XBEInfo struct {
	Timestamp int64
	Checksum  int64
	Name      string
}

// NewXBEInfo ..
func NewXBEInfo(body []string) *XBEInfo {
	merged := strings.Join(body, " ")
	data := helpers.ParseSpaceSeparatedValues(merged)

	return &XBEInfo{
		Timestamp: helpers.ConvertHexToInt64(data["timestamp"]),
		Checksum:  helpers.ConvertHexToInt64(data["checksum"]),
		Name:      data["name"],
	}
}
