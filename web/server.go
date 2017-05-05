package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	models "../models"
	proIntegrator "../proto"
	"github.com/gorilla/mux"
)

//Server encapsulates web interface features
type Server struct {
	ComponentHandlers *[]models.ComponentHandler
	AddChain          models.AddChain
	router            *mux.Router
}

//InIt method initialzes web interface
func (srv *Server) InIt() {
	srv.router = mux.NewRouter()
	srv.router.HandleFunc("/components", srv.getComponents).Methods("GET")
	srv.router.HandleFunc("/chain", srv.addChain).Methods("POST")
	srv.router.PathPrefix("/web").Handler(http.StripPrefix("/web", http.FileServer(http.Dir("web/console"))))

	go http.ListenAndServe(":8020", srv.router)
}

//UpdateInterface method add new endpoints for taking entry from external
func (srv *Server) UpdateInterface(url string, handler models.WebHandler, reqType string) {
	srv.router.HandleFunc(url, handler).Methods(reqType)
}

func (srv *Server) getComponents(w http.ResponseWriter, req *http.Request) {
	components := toComponents(srv.ComponentHandlers)
	json.NewEncoder(w).Encode(components)
}

func (srv *Server) addChain(w http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)
	srv.AddChain(proIntegrator.SaveChainRequest{b})
}

func toComponents(compHandlers *[]models.ComponentHandler) []*proIntegrator.Component {
	components := make([]*proIntegrator.Component, 0)

	for _, compHandler := range *compHandlers {
		components = append(components, compHandler.Component)
	}
	return components
}
