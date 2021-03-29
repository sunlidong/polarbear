package main

import (
	"pol/server"
	"pol/router"
)

func main() {
	server.Main_cryptogen()

	// router

	router.GetAllRounters()
}
