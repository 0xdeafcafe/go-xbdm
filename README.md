go-xbdm
===
[![Build Status](https://img.shields.io/travis/0xdeafcafe/go-xbdm/master.svg?style=flat-square&label=build%20status)](https://travis-ci.org/0xdeafcafe/go-xbdm)
[![Github Issues](https://img.shields.io/github/issues/0xdeafcafe/go-xbdm.svg?style=flat-square)](https://github.com/0xdeafcafe/go-xbdm/issues)
[![Github License](https://img.shields.io/github/license/0xdeafcafe/go-xbdm.svg?style=flat-square)](https://github.com/0xdeafcafe/go-xbdm/blob/master/LICENSE)

A go library for interacting with an Xbox 360 Development Kit.

## Installation
```
go get github.com/0xdeafcafe/go-xbdm
```

## Usage
```
import "github.com/0xdeafcafe/go-xbdm"
````

This example with be based on the assumption that the Xbox 360 Developer Kit's IP
address is `192.168.1.88`.

``` go
func main() {
  xbdm, err := goxbdm.NewXBDMClient("192.168.1.88")
  if err != nil {
    fmt.Println("There was an error connecting to the Development Kit.")
    return
  }

  // Set some random memory offset
  resp, err := xbdm.SetMemory(0xBF9BC9CC, "0000fa67")
  if err != nil {
    fmt.Println("There was an error writing memory to the Development Kit.")
    return
  }

  fmt.Println(resp)
}
```
