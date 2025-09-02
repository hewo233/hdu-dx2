package main

import (
	"github.com/hewo233/hdu-dx2/initall"
	"github.com/hewo233/hdu-dx2/route"
	"log"
)

func main() {
	initall.Init()

	route.InitRoute()

	err := route.R.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
