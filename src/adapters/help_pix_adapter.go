package adapters

import (
	"laches1sm/help_pix_go/infrastructure"
	"log"
)

type HelpPixHTTPAdapter struct {
	Logger log.Logger
	Infra  infrastructure.HelpPixInfra
}

func NewHelpPixAdapter(logger log.Logger, infra infrastructure.HelpPixInfra) *HelpPixHTTPAdapter {
	return &HelpPixHTTPAdapter{
		Logger: logger,
		Infra:  infra,
	}
}
