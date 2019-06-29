package main

import (
	"log"

	wtmfactor "github.com/prongbang/wtm-factor"
)

func main() {

	wtm := wtmfactor.NewWtmFactor(wtmfactor.WtmConfig{
		URL:       "https://myweb.com",
		FactorRow: ".form-row .py-1",
		FactorCol: ".card",
	})

	algorithms := wtm.GetFactor()

	log.Println("algorithms -> ", algorithms)
}
