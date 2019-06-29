# WTM Algorithm Factory

### Install

```bash
go get -u github.com/prongbang/wtm-factor
```

### How to use

```golang
package main

import (
	"log"

	wtmfactor "github.com/prongbang/wtm-factor"
)

func main() {

	wtm := wtmfactor.NewWtmFactor(wtmfactor.WtmConfig{
		URL: "https://myweb.com",
	})

	algorithms := wtm.GetFactor()

	log.Println("algorithms -> ", algorithms)
}

```
