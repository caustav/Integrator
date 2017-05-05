package chain

import (
	"strconv"
	"strings"
	"sync"

	"fmt"

	models "../models"
	"github.com/Jeffail/gabs"
)

//Runner excutes the one chain per request
type Runner struct {
	ComponentHandlers []models.ComponentHandler
	ComponentName     string
	Containers        []Container
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
		if source.IsInExecution == false {
			wg.Add(1)
			go source.execute(wg)
		}
		return
	}
	sourceOutputParam, _ := strconv.Unquote(jsonParsed.Path(comName + ".src").Index(counter).String())
	targetstr, _ := strconv.Unquote(jsonParsed.Path(comName + ".dest").Index(counter).String())
	stringSlice := strings.Split(targetstr, ".")
	targetName := stringSlice[0]
	targetParam := stringSlice[2]
	if target == nil {
		target = chainRunner.GetContainer(targetName)
	}
	source.AddParam(OperatingParam{target, targetParam, sourceOutputParam})
	counter++
	chainRunner.runRecursive(source, target, jsonParsed, counter, wg)
	if counter == 1 {
		chainRunner.runRecursive(target, nil, jsonParsed, 0, wg)
	}
}

//Init initialize object
func (chainRunner *Runner) Init(compHandlers []models.ComponentHandler) {
	chainRunner.ComponentHandlers = make([]models.ComponentHandler, 0)
	chainRunner.Containers = make([]Container, 0)
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
		container := Container{compHlr, make(chan bool), false, make([]OperatingParam, 0)}
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

	componentHandler := *container.ComponentHandler
	component := *container.ComponentHandler.Component
	componentHandler.Component = &component
	con := &Container{&componentHandler, make(chan bool), false, make([]OperatingParam, 0)}
	return con
}
