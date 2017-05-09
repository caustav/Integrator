package components

import (
	"fmt"

	"golang.org/x/net/context"

	proIntegrator "../proto"
)

// SampleIn is used for taking input from external source, it contains a component object
type SampleIn struct {
	Component proIntegrator.Component
}

// FetchComponent method reveals deatils about the component it contains
func (in *SampleIn) FetchComponent() *proIntegrator.Component {
	in.Component.Name = "SampleIn"

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
func (in *SampleIn) Execute(ctx context.Context, req *proIntegrator.ExecuteRequest) (*proIntegrator.ExecuteResponse, error) {

	fmt.Println("SampleIn::Execute called with " + req.Component.ParamsInput["Param1"].Str + " " +
		req.Component.ParamsInput["Param2"].Str)

	in.Component.ParamsOutput["Param1"].Str = req.Component.ParamsInput["Param1"].Str + " SampleIn"
	in.Component.ParamsOutput["Param2"].Str = req.Component.ParamsInput["Param2"].Str + " SampleIn"

	component := &proIntegrator.Component{}
	mapIn := make(map[string]*proIntegrator.DataType)
	mapOut := make(map[string]*proIntegrator.DataType)
	for key, value := range in.Component.ParamsInput {
		mapIn[key] = value
	}
	for key, value := range in.Component.ParamsOutput {
		mapOut[key] = value
	}
	component.Name = in.Component.Name
	component.ParamsInput = mapIn
	component.ParamsOutput = mapOut

	response := &proIntegrator.ExecuteResponse{component}

	return response, nil
}
