package chain

import (
	"strconv"
	"strings"
	"sync"

	"fmt"

	models "../models"
	proIntegrator "../proto"
	"github.com/Jeffail/gabs"
)

//Runner excutes the one chain per request
type Runner struct {
	ComponentHandlers []models.ComponentHandler
	ComponentName     string
	Containers        []Container
	ContainerMap      map[string]*Container
}

//Run executes a chain
func (chainRunner *Runner) Run(bytes []byte) {
	jsonParsed, _ := gabs.ParseJSON(bytes)
	comName := chainRunner.ComponentName
	wg := &sync.WaitGroup{}
	chainRunner.runRecursive(chainRunner.GetContainer(comName), nil, jsonParsed, 0, wg)
	wg.Wait()
	fmt.Println("Run ends here ...")
}

func (chainRunner *Runner) runRecursive(source *Container, target *Container, jsonParsed *gabs.Container,
	counter int, wg *sync.WaitGroup) {

	comName := source.ComponentHandler.Component.Name

	if !jsonParsed.Path(comName + ".src").Index(counter).Exists() {
		if source.IsInExecution == 0 {
			wg.Add(1)
			go source.execute(wg)
			source.IsInExecution = 3
		}
		return
	}
	sourceOutputParam, _ := strconv.Unquote(jsonParsed.Path(comName + ".src").Index(counter).String())
	targetstr, _ := strconv.Unquote(jsonParsed.Path(comName + ".dest").Index(counter).String())
	stringSlice := strings.Split(targetstr, ".")
	targetName := stringSlice[0]
	targetParam := stringSlice[2]
	targetID := stringSlice[3]
	target = chainRunner.getTarget(targetName, targetID)
	source.AddParam(OperatingParam{target, targetParam, sourceOutputParam})
	counter++
	chainRunner.runRecursive(source, target, jsonParsed, counter, wg)
	if target.IsInExecution == 0 {
		chainRunner.runRecursive(target, nil, jsonParsed, 0, wg)
	}
}

func (chainRunner *Runner) getTarget(targetName string, id string) *Container {
	val := chainRunner.ContainerMap[id]
	if val == nil {
		val = chainRunner.GetContainer(targetName)
		chainRunner.ContainerMap[id] = val
	}
	return val
}

//Init initialize object
func (chainRunner *Runner) Init(compHandlers []models.ComponentHandler) {
	chainRunner.ComponentHandlers = make([]models.ComponentHandler, 0)
	chainRunner.Containers = make([]Container, 0)
	chainRunner.ContainerMap = make(map[string]*Container)
	for _, compHandler := range compHandlers {
		compHlr := &models.ComponentHandler{nil, nil, nil}

		if compHandler.ClientConn != nil {
			cc := *compHandler.ClientConn
			compHlr.ClientConn = &cc
		}

		if compHandler.Component != nil {
			co := *compHandler.Component
			compHlr.Component = &co
		}

		compHlr.CompService = compHandler.CompService
		chainRunner.ComponentHandlers = append(chainRunner.ComponentHandlers, *compHlr)
		container := Container{compHlr, make(chan bool), 0, make([]OperatingParam, 0)}
		chainRunner.Containers = append(chainRunner.Containers, container)
	}
}

//GetComponentHandler to fetch component
func (chainRunner *Runner) GetComponentHandler(name string) *models.ComponentHandler {
	var compHandler *models.ComponentHandler
	retValue := false
	for _, componentHandler := range chainRunner.ComponentHandlers {
		compHandler = &componentHandler
		if compHandler.Component.Name == name {
			retValue = true
			break
		}
	}
	if retValue == false {
		compHandler = nil
	}
	return compHandler
}

//GetContainer gives container from component name
func (chainRunner *Runner) GetContainer(name string) *Container {
	var container Container
	for _, val := range chainRunner.Containers {
		if name == val.ComponentHandler.Component.Name {
			container = val
			break
		}
	}
	componentHandler := &models.ComponentHandler{}
	componentHandler.ClientConn = container.ComponentHandler.ClientConn
	componentHandler.CompService = container.ComponentHandler.CompService

	component := &proIntegrator.Component{}
	mapIn := make(map[string]*proIntegrator.DataType)
	mapOut := make(map[string]*proIntegrator.DataType)
	for key, value := range container.ComponentHandler.Component.ParamsInput {
		mapIn[key] = value
	}
	for key := range container.ComponentHandler.Component.ParamsOutput {
		mapOut[key] = nil
	}
	component.Name = container.ComponentHandler.Component.Name
	component.ParamsInput = mapIn
	component.ParamsOutput = mapOut

	componentHandler.Component = component

	con := &Container{componentHandler, make(chan bool), 0, make([]OperatingParam, 0)}
	return con
}
