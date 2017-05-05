package models

import (
	"net/http"

	"golang.org/x/net/context"

	proIntegrator "../proto"
	"google.golang.org/grpc"
)

//ComponentHandler associates proto component and client connection
type ComponentHandler struct {
	ClientConn  *grpc.ClientConn
	Component   *proIntegrator.Component
	CompService ComponentService
}

// AddChain provides callback for add chain features to service from server
type AddChain func(req proIntegrator.SaveChainRequest) bool

//WebHandler provides callback for to handle request and response for server
type WebHandler func(w http.ResponseWriter, req *http.Request)

//WebInterface encapsultres neccessary components
type WebInterface struct {
	URL                string
	WebHandlerFunction WebHandler
	MethodType         string
}

//MapChainInfo maps chain info with chain id
type MapChainInfo map[string][]byte

//IService is basic type of integrator service
type IService interface {
	Start(ctx context.Context, request *proIntegrator.StartRequest) (*proIntegrator.StartResponse, error)
}

//ComponentService is only for components owned by integrator
type ComponentService interface {
	Execute(ctx context.Context, req *proIntegrator.ExecuteRequest) (*proIntegrator.ExecuteResponse, error)
}
