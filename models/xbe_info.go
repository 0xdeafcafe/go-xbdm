package models

// XBEInfo defines the structure of the `xbeinfo running` response
type XBEInfo struct {
	Timestamp int64  `goxbdm:"timestamp"`
	Checksum  int64  `goxbdm:"checksum"`
	Name      string `goxbdm:"name"`
}
