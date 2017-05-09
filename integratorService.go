package main

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Jeffail/gabs"

	"golang.org/x/net/context"

	chain "./chain"
	input "./components"
	models "./models"
	proIntegrator "./proto"
	webConsole "./web"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type integratorService struct {
	ComponentHandlers []models.ComponentHandler
	MapChain          models.MapChainInfo
}

func main() {
	var intgtrService integratorService
	intgtrService.InIt()
	var wg sync.WaitGroup
	wg.Add(2)
	go listen(&intgtrService, ":3020", &wg)
	go setupConsole(&intgtrService, &wg)
	wg.Wait()
}

func listen(intgtrService *integratorService, port string, wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srvr := grpc.NewServer()
	proIntegrator.RegisterIntegratorServer(srvr, intgtrService)
	reflection.Register(srvr)
	if err := srvr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	wg.Done()
}

func setupConsole(intgtrService *integratorService, wg *sync.WaitGroup) {
	var wConsole webConsole.Server
	wConsole.ComponentHandlers = &intgtrService.ComponentHandlers
	wConsole.InIt()
	wConsole.AddChain = intgtrService.AddChain
	intgtrService.initComponents(&wConsole)
	wg.Done()
}

func (intgService *integratorService) initComponents(webConsole *webConsole.Server) {
	var in input.Input
	in.Service = intgService
	componentHandler := models.ComponentHandler{nil, in.FetchComponent(), &in}
	intgService.ComponentHandlers = append(intgService.ComponentHandlers, componentHandler)

	var sampleIn input.SampleIn
	componentHandler = models.ComponentHandler{nil, sampleIn.FetchComponent(), &sampleIn}
	intgService.ComponentHandlers = append(intgService.ComponentHandlers, componentHandler)

	var sampleIn1 input.SampleIn1
	componentHandler = models.ComponentHandler{nil, sampleIn1.FetchComponent(), &sampleIn1}
	intgService.ComponentHandlers = append(intgService.ComponentHandlers, componentHandler)

	var fileOut input.FileOut
	componentHandler = models.ComponentHandler{nil, fileOut.FetchComponent(), &fileOut}
	intgService.ComponentHandlers = append(intgService.ComponentHandlers, componentHandler)

	webInterface := in.FetchWebInterface()
	webConsole.UpdateInterface(webInterface.URL, webInterface.WebHandlerFunction, webInterface.MethodType)
}

func (intgService *integratorService) InIt() {
	intgService.ComponentHandlers = make([]models.ComponentHandler, 0)
	intgService.MapChain = make(models.MapChainInfo)
}

func (intgService *integratorService) Register(ctx context.Context, module *proIntegrator.Module) (*proIntegrator.RegisterResponse, error) {
	retVal := &proIntegrator.RegisterResponse{true}
	conn, err := grpc.Dial(module.Url, grpc.WithInsecure())
	intgService.PopulateComponentHandlers(conn, module)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return retVal, err
}

func (intgService *integratorService) Start(ctx context.Context, request *proIntegrator.StartRequest) (*proIntegrator.StartResponse, error) {
	retVal := &proIntegrator.StartResponse{true}
	chainInfo := intgService.MapChain[request.ChainId]
	component := request.StartComponent

	compHandlerStart := intgService.GetComponentHandler(component.Name)
	exReq := &proIntegrator.ExecuteRequest{component}
	if compHandlerStart.ClientConn != nil {

	} else {
		compHandlerStart.CompService.Execute(ctx, exReq)
	}
	intgService.startChain(chainInfo, component.Name)

	return retVal, nil
}

func (intgService *integratorService) startChain(bytes []byte, comName string) {
	runner := &chain.Runner{nil, comName, nil}
	runner.Init(intgService.ComponentHandlers)
	runner.Run(bytes)
}

func (intgService *integratorService) SaveChain(ctx context.Context, in *proIntegrator.SaveChainRequest) (*proIntegrator.SaveChainResponse, error) {
	retVal := &proIntegrator.SaveChainResponse{true}
	if in == nil {
		log.Fatalf("SaveChainRequest is null")
	}
	jsonParsed, _ := gabs.ParseJSON(in.ChainInfo)
	floatID := jsonParsed.Path("ChainId").Data().(float64)
	strChainID := strconv.FormatFloat(floatID, 'f', -1, 64)
	intgService.MapChain[strChainID] = in.ChainInfo
	return retVal, nil
}

func (intgService *integratorService) AddChain(req proIntegrator.SaveChainRequest) bool {
	jsonParsed, _ := gabs.ParseJSON(req.ChainInfo)
	floatID := jsonParsed.Path("ChainId").Data().(float64)
	strChainID := strconv.FormatFloat(floatID, 'f', -1, 64)
	intgService.MapChain[strChainID] = req.ChainInfo
	return true
}

func (intgService *integratorService) PopulateComponentHandlers(conn *grpc.ClientConn, module *proIntegrator.Module) {
	components := module.Components
	for _, component := range components {
		componentHandler := models.ComponentHandler{conn, component, nil}
		intgService.ComponentHandlers = append(intgService.ComponentHandlers, componentHandler)
	}
}

func (intgService *integratorService) GetComponentHandler(name string) *models.ComponentHandler {
	compHandler := &models.ComponentHandler{}
	retValue := false
	for _, componentHandler := range intgService.ComponentHandlers {
		if componentHandler.Component.Name == name {
			if componentHandler.ClientConn != nil {
				clientConn := *componentHandler.ClientConn
				compHandler.ClientConn = &clientConn
			}
			component := *componentHandler.Component
			compHandler.Component = &component
			compHandler.CompService = componentHandler.CompService
			retValue = true
			break
		}
	}
	if retValue == false {
		compHandler = nil
	}
	return compHandler
}
