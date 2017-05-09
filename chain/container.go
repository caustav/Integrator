package chain

import (
	"context"
	"sync"

	models "../models"
	proIntegrator "../proto"
)

//Container holds Component
type Container struct {
	ComponentHandler *models.ComponentHandler
	Blocked          chan bool
	IsInExecution    bool
	OperatingParams  []OperatingParam
}

//OperatingParam contains parameters to be used in execution
type OperatingParam struct {
	Target         *Container
	TargetInParam  string
	SourceOutParam string
}

//InIt intilaizes the object
func (container *Container) InIt(comp *proIntegrator.Component) {
	container.ComponentHandler.Component = comp
}

//SetComponent sets value in the component inside
func (container *Container) SetComponent(k string, v *proIntegrator.DataType) {
	container.ComponentHandler.Component.ParamsInput[k] = v
	container.manageInputs()
}

//AddParam adds element to the collection
func (container *Container) AddParam(param OperatingParam) {
	container.OperatingParams = append(container.OperatingParams, param)
}

//Execute component inside
func (container *Container) execute(wg *sync.WaitGroup) {
	if container.isAllInputFilled() == false {
		container.IsInExecution = true
		blocked := <-container.Blocked
		if blocked == false {
			container.executeComponent()
		}
		container.IsInExecution = false
	} else {
		container.executeComponent()
	}

	for i := 0; i < len(container.OperatingParams); i++ {
		// r.setParam(r.params[i].obj, r.params[i].targetParam, r.params[i].outVal)
		container.OperatingParams[i].Target.SetComponent(container.OperatingParams[i].TargetInParam,
			container.ComponentHandler.Component.ParamsOutput[container.OperatingParams[i].SourceOutParam])
	}

	wg.Done()
}

func (container *Container) executeComponent() {
	if container.ComponentHandler.ClientConn != nil {
		moduleClient := proIntegrator.NewIntegratorModuleClient(container.ComponentHandler.ClientConn)
		response, _ := moduleClient.Execute(context.Background(), &proIntegrator.ExecuteRequest{container.ComponentHandler.Component})
		container.ComponentHandler.Component = response.Component
	} else {
		request := &proIntegrator.ExecuteRequest{container.ComponentHandler.Component}
		container.ComponentHandler.CompService.Execute(context.Background(), request)
	}
}

func (container *Container) manageInputs() {
	if container.isAllInputFilled() == true && container.IsInExecution == true {
		container.Blocked <- false
	}
}

func (container *Container) isAllInputFilled() bool {
	ret := true
	for _, value := range container.ComponentHandler.Component.ParamsInput {
		if value == nil || len(value.Str) <= 0 {
			ret = false
			break
		}
	}
	return ret
}
