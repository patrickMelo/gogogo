package lib

import (
	"os"
	"os/signal"
	"runtime"
	"strings"
)

// Starts the service infrastructure.
func Start(name string, versionString string) {
	serviceName = name
	serviceVersionString = versionString

	LogInformation(serviceLogTag, "%s - version %s", serviceName, serviceVersionString)

	var httpOptions = GenericMap{
		"http.listenAddress": ":8080",
		"http.keepAlive":     true,
	}

	if !startHttp(httpOptions) {
		return
	}

	LogInformation(serviceLogTag, "started")
}

// Runs the service (start handling requests).
func Run() {
	LogVerbose(serviceLogTag, "running")

	go waitSignal()

	serviceRunning = true
	for serviceRunning {
		runtime.Gosched()
	}

	close(serviceSignalChannel)
	stopHttp()

	LogInformation(serviceLogTag, "stopped")
}

func handleRequest(request *Request, response *Response) error {
	LogInformation(serviceLogTag, "(%s) %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	LogVerbose(serviceLogTag, "(%s) searching route for %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	var route = findRoute(request.Type, request.Path)

	if route == nil {
		response.Status = ResourceNotFoundStatus
		return nil
	}

	LogVerbose(serviceLogTag, "(%s) route found for %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	if !route.isPublic && (request.AuthToken == "") {
		LogVerbose(serviceLogTag, "(%s) route %s:%s is not public and no authorization token was specified", request.Id, strings.ToLower(request.Type.String()), request.Path)
		response.Status = AuthenticationRequiredStatus
		return nil
	}

	if route.contract != nil {
		var contractErrors = route.contract.Validate(request.Data)

		if len(contractErrors) > 0 {
			LogVerbose(serviceLogTag, "(%s) payload has contract validation errors: %v", request.Id, contractErrors)

			response.Status = InvalidDataStatus
			response.Data[ContractValidationErrorsFieldName] = contractErrors

			return nil
		}
	}

	LogVerbose(serviceLogTag, "(%s) handling request", request.Id)

	return route.handler(request, response)
}

var (
	serviceRunning       bool = false
	serviceName          string
	serviceVersionString string
	serviceSignalChannel chan os.Signal
)

const serviceLogTag = "service"

func waitSignal() {
	serviceSignalChannel = make(chan os.Signal, 1)
	signal.Notify(serviceSignalChannel, os.Interrupt)
	<-serviceSignalChannel
	serviceRunning = false
}
