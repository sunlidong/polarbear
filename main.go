package main

import (
	"des/desgin"
)

func main() {
	lishi := desgin.NewLiShi()
	yuwen := desgin.NewYuWen()

	// buffer
	lishi.Set("历史")
	yuwen.Set("语文")
	lishi.Say("lishi")
	yuwen.Say("语文")

	api := desgin.NewAPI()
	ret := api.Test()
}
