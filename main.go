package main
import(
"pol/cmd/cryptogen"
"pol/router"

)

func main() {
	cryptogen.Main_cryptogen()

	// router

	router.GetAllRounters()
}
