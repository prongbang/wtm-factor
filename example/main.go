package main

import (
	"log"

	wtmfactor "github.com/prongbang/wtm-factor"
)

func main() {

	wtm := wtmfactor.NewWtmFactor(wtmfactor.WtmConfig{
		URL: "https://whattomine.com/",
	})

	algorithms := wtm.GetFactorName()
	algorithmKey := wtm.GetFactorKey()

	log.Println("algorithms -> ", algorithms)
	log.Println("algorithmKey -> ", algorithmKey)

}
