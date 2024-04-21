package main

import (
	"github.com/mimminou/UrlCleaner-ByFood/server"
)

func main() {
	var PORT uint16 = 50456
	server.Serve(PORT)
}
