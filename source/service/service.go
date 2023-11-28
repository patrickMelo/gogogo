package service

import (
	"gogogo/log"
	"gogogo/requests"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

type Listener interface {
	Start() error
	Stop()
}

// Initializes the service infrastructure.
func Start(name string, versionString string, listeners ...Listener) {
	log.Information(_logTag, "%s - version %s", name, versionString)
	_listeners = listeners
}

// Starts listening for requests.
func Run() (err error) {
	defer log.Information(_logTag, "stopped")

	for _, listener := range _listeners {
		if err = listener.Start(); err != nil {
			stopListeners()
			return
		}
	}

	log.Information(_logTag, "running")

	go waitSignal()

	_isRunning = true

	for _isRunning {
		runtime.Gosched()
	}

	close(_signalChannel)
	stopListeners()

	return
}

// Stops the service.
func Stop() {
	_isRunning = false
}

func HandleRequest(request *requests.Request, response *requests.Response) error {
	log.Information(_logTag, "(%s) %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	log.Verbose(_logTag, "(%s) searching route for %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	var route = findRoute(request.Type, request.Path)

	if route == nil {
		response.Status = requests.ResourceNotFound
		return nil
	}

	log.Verbose(_logTag, "(%s) route found for %s:%s", request.Id, strings.ToLower(request.Type.String()), request.Path)

	if !route.isPublic && (request.Metadata.GetString("Token", "") == "") {
		log.Verbose(_logTag, "(%s) route %s:%s is not public and no authorization token was specified", request.Id, strings.ToLower(request.Type.String()), request.Path)
		response.Status = requests.AuthenticationRequired
		return nil
	}

	if route.contract != nil {
		var contractErrors = route.contract.Validate(request.Data)

		if len(contractErrors) > 0 {
			log.Verbose(_logTag, "(%s) payload has contract validation errors: %v", request.Id, contractErrors)

			response.Status = requests.InvalidData
			response.Data["errors"] = contractErrors

			return nil
		}
	}

	log.Verbose(_logTag, "(%s) handling request", request.Id)

	return route.handler(request, response)
}

var (
	_listeners     []Listener
	_isRunning     bool
	_signalChannel chan os.Signal
)

const (
	_logTag = "service"
)

func waitSignal() {
	_signalChannel = make(chan os.Signal, 1)
	signal.Notify(_signalChannel, os.Interrupt)
	<-_signalChannel
	Stop()
}

func stopListeners() {
	for _, listener := range _listeners {
		listener.Stop()
	}
}
