package components

import (
	"fmt"

	"golang.org/x/net/context"

	proIntegrator "../proto"
)

// FileOut is used for taking input from external source, it contains a component object
type FileOut struct {
	Component proIntegrator.Component
}

// FetchComponent method reveals deatils about the component it contains
func (in *FileOut) FetchComponent() *proIntegrator.Component {
	in.Component.Name = "FileOut"

	mapParamIn := make(map[string]*proIntegrator.DataType)
	mapParamOut := make(map[string]*proIntegrator.DataType)

	mapParamIn["data"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}

	mapParamOut["File"] = &proIntegrator.DataType{Type: proIntegrator.DataType_STR}

	in.Component.ParamsInput = mapParamIn
	in.Component.ParamsOutput = mapParamOut
	return &in.Component
}

// Execute method runs the core functionality it contains
func (in *FileOut) Execute(ctx context.Context, req *proIntegrator.ExecuteRequest) (*proIntegrator.ExecuteResponse, error) {
	fmt.Println(in.Component.Name + " is called with " + req.Component.ParamsInput["data"].Str)
	component := &proIntegrator.Component{}
	mapIn := make(map[string]*proIntegrator.DataType)
	mapOut := make(map[string]*proIntegrator.DataType)
	for key, value := range req.Component.ParamsInput {
		mapIn[key] = value
	}
	mapOut["File"] = req.Component.ParamsInput["data"]
	mapOut["File"].Str += " FileOut added"

	component.Name = in.Component.Name
	component.ParamsInput = mapIn
	component.ParamsOutput = mapOut

	response := &proIntegrator.ExecuteResponse{component}
	return response, nil
}
