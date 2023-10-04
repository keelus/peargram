package main

import (
	"peargram/server"
	"peargram/utils"
)

func main() {
	utils.LoadEnv()

	server.Init()
}
