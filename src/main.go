package main

import (
	"log"

	"github.com/laches1sm/help_pix_go/src/adapters"
	"github.com/laches1sm/help_pix_go/src/httpserver"
	"github.com/laches1sm/help_pix_go/src/infrastructure"
)

func main() {
	log := log.New(log.Writer(), "help-pix", 0)
    var infra infrastructure.HelpPixInfra
	adapter := adapters.NewHelpPixAdapter(log, infra)

	server := httpserver.NewParrotServer(log, adapter)
	server.SetupRoutes()
	if err := server.Start(httpserver.ServerPort); err != nil {
		log.Println(err.Error())
	}

}
