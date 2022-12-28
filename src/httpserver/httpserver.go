package httpserver

import (
	"github.com/laches1sm/help_pix_go/src/adapters"
	"log"
	"net/http"
)

const (
	ServerPort       = ":7000"
	UserDataEndpoint = "/send-user-data"
	ResultEndpoint   = "/result"
)

// HelpPixServer is an interface to an HTTP server which handles requests
type HelpPixServer struct {
	Mux                *http.ServeMux
	Logger             log.Logger
	HelpPixHTTPAdapter *adapters.HelpPixHTTPAdapter
}

// SetupRoutes configures the routes of the API
func (srv *HelpPixServer) SetupRoutes() {
	srv.Mux.Handle(UserDataEndpoint, http.HandlerFunc(srv.HelpPixHTTPAdapter.GetSummonerInfo))
	// srv.Mux.Handle(ResultEndpoint, http.HandlerFunc(srv.ParrotHTTPAdapter.AddParrot))
}

// Start sets up the HTTP webserver to listen and handle traffic. It
// takes the port number to listen on as a parameter in the form ":PORT_NUMBER"
func (srv *HelpPixServer) Start(port string) error {
	return http.ListenAndServe(port, srv.Mux)
}

// NewParrotServer returns an instance of a configured ParrotServer
func NewParrotServer(logger log.Logger, adapter *adapters.HelpPixHTTPAdapter) *HelpPixServer {
	httpServer := &HelpPixServer{
		Mux:                http.NewServeMux(),
		Logger:             logger,
		HelpPixHTTPAdapter: adapter,
	}
	return httpServer
}
