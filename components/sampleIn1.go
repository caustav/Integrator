package components

import (
	"fmt"

	"golang.org/x/net/context"

	proIntegrator "../proto"
)

// SampleIn1 is used for taking input from external source, it contains a component object
type SampleIn1 struct {
	Component proIntegrator.Component
}

// FetchComponent method reveals deatils about the component it contains
func (in *SampleIn1) FetchComponent() *proIntegrator.Component {
	in.Component.Name = "SampleIn1"

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
func (in *SampleIn1) Execute(ctx context.Context, req *proIntegrator.ExecuteRequest) (*proIntegrator.ExecuteResponse, error) {
	response := &proIntegrator.ExecuteResponse{}
	fmt.Println("SampleIn1::Execute called with " + in.Component.ParamsInput["Param1"].Str + " " +
		in.Component.ParamsInput["Param2"].Str)
	in.Component.ParamsOutput["Param1"].Str = in.Component.ParamsInput["Param1"].Str + " SampleIn1"
	in.Component.ParamsOutput["Param2"].Str = in.Component.ParamsInput["Param2"].Str + " SampleIn1"
	return response, nil
}
