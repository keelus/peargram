package server

import (
	"runtime"
)

func Init() {
	router := SetupRouter()

	port := ":9990"
	if runtime.GOOS == "windows" {
		port = "localhost:80"
	}

	router.Run(port)
}
