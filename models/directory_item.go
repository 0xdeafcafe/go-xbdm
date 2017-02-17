package models

import "fmt"

// DirectoryItem defines the structure of the `dirlist name=":dir"` response
type DirectoryItem struct {
	Name string

	SizeHi int64
	SizeLo int64

	CreateHi int64
	CreateLo int64

	ChangeHi int64
	ChangeLo int64

	IsDirectory bool
}

// NewDirectoryItem ..
func NewDirectoryItem(body string) *DirectoryItem {
	format := `name=%s sizehi=0x%x sizelo=0x%x createhi=0x%x createlo=0x%x changehi=0x%x changelo=0x%x %s`
	var name, directory string
	var sizeHi, sizeLo, createHi, createLo, changeHi, changeLo int64
	fmt.Sscanf(body, format, &name, &sizeHi, &sizeLo, &createHi, &createLo, &changeHi, &changeLo, &directory)
	isDirectory := directory == "directory"

	return &DirectoryItem{
		Name: name,

		SizeHi: sizeHi,
		SizeLo: sizeLo,

		CreateHi: createHi,
		CreateLo: createLo,

		ChangeHi: changeHi,
		ChangeLo: changeLo,

		IsDirectory: isDirectory,
	}
}
