package main

import (
	"log"

	"github.com/laches1sm/help_pix_go/src/adapters"
	"github.com/laches1sm/help_pix_go/src/httpserver"
)

func main() {
	log := log.New(log.Writer(), "help-pix-", 0)
	adapter := adapters.NewHelpPixAdapter(log)

	server := httpserver.NewParrotServer(log, adapter)
	log.Print(`Creating server`)
	server.SetupRoutes()
	if err := server.Start(httpserver.ServerPort); err != nil {
		log.Println(err.Error())
	}
}
