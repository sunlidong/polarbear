package main

import (
	"pol/router"
	"pol/server"
)

func main() {

	var i server.Generate
	i.Generate()

	// router

	router.GetAllRounters()
}
