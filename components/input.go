package components

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/Jeffail/gabs"

	"fmt"

	models "../models"
	proIntegrator "../proto"
)

// Input is used for taking input from external source, it contains a component object
type Input struct {
	Component proIntegrator.Component
	Service   models.IService
}

type inputDataParam struct {
	Param1 string
	Param2 string
}

// FetchComponent method reveals deatils about the component it contains
func (in *Input) FetchComponent() *proIntegrator.Component {
	in.Component.Name = "Input"

	mapParamIn := make(map[string]*proIntegrator.DataType)
	mapParamOut := make(map[string]*proIntegrator.DataType)

	mapParamIn["Param1"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}
	mapParamIn["Param2"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}

	mapParamOut["Param1"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}
	mapParamOut["Param2"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}

	in.Component.ParamsInput = mapParamIn
	in.Component.ParamsOutput = mapParamOut
	return &in.Component
}

// Execute method runs the core functionality it contains
func (in *Input) Execute(ctx context.Context, req *proIntegrator.ExecuteRequest) (*proIntegrator.ExecuteResponse, error) {
	response := &proIntegrator.ExecuteResponse{}
	component := req.Component
	component.ParamsOutput["Param1"] = component.ParamsInput["Param1"]
	component.ParamsOutput["Param2"] = component.ParamsInput["Param2"]
	return response, nil
}

// FetchWebInterface method provides endpoint information
func (in *Input) FetchWebInterface() models.WebInterface {

	var webInterface models.WebInterface
	webInterface.URL = "/input"
	webInterface.WebHandlerFunction = in.EndpointHandler
	webInterface.MethodType = "POST"
	return webInterface
}

// EndpointHandler method takes input from consuming component
func (in *Input) EndpointHandler(w http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)
	jsonParsed, _ := gabs.ParseJSON(b)
	chainID := jsonParsed.Path("ChainId").Data().(float64)
	param1 := jsonParsed.Path("Param1").Data().(string)
	param2 := jsonParsed.Path("Param2").Data().(string)

	in.Component.ParamsInput["Param1"].Str = param1
	in.Component.ParamsInput["Param2"].Str = param2

	var request proIntegrator.StartRequest
	request.ChainId = strconv.FormatFloat(chainID, 'f', -1, 64)
	request.StartComponent = &in.Component
	in.Service.Start(nil, &request)
	fmt.Println("Finished ...")
}
