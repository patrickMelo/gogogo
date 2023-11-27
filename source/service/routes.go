package service

import (
	"fmt"
	"gogogo/data"
	"gogogo/log"
	"gogogo/requests"
	"strings"
)

func AddPublicPull(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Pull, path, true, handler, contract)
}

func AddPublicPush(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Push, path, true, handler, contract)

}

func AddPublicUpdate(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Update, path, true, handler, contract)

}

func AddPublicDelete(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Delete, path, true, handler, contract)

}

func AddPrivatePull(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Pull, path, false, handler, contract)

}

func AddPrivatePush(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Push, path, false, handler, contract)
}

func AddPrivateUpdate(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Update, path, false, handler, contract)
}

func AddPrivateDelete(path string, handler requests.Handler, contract *data.Contract) {
	addRoute(requests.Delete, path, false, handler, contract)
}

type routePartInfo struct {
	isVariable bool
	part       string
}

type routeInfo struct {
	requestType requests.Type
	isPublic    bool
	handler     requests.Handler
	parts       []routePartInfo
	contract    *data.Contract
}

var (
	_routes = make(map[string]*routeInfo)
)

func findRoute(requestType requests.Type, path string) *routeInfo {
	var rawPathParts = strings.Split(path, "/")
	var pathParts = make([]string, 0, len(rawPathParts))

	for _, part := range rawPathParts {
		if part != "" {
			pathParts = append(pathParts, part)
		}
	}

	if len(pathParts) == 0 {
		return nil
	}

	for _, route := range _routes {
		if (route.requestType != requestType) || (len(route.parts) != len(pathParts)) {
			continue
		}

		var routeFound = true

		for index, routePart := range route.parts {
			if routePart.isVariable {
				continue
			}

			if routePart.part != pathParts[index] {
				routeFound = false
				break
			}
		}

		if routeFound {
			return route
		}
	}

	return nil
}

func addRoute(requestType requests.Type, path string, isPublic bool, handler requests.Handler, contract *data.Contract) {
	log.Verbose(_logTag, "adding route for %s:%s", strings.ToLower(requestType.String()), path)

	if existingRoute := findRoute(requestType, path); existingRoute != nil {
		log.Error(_logTag, fmt.Errorf("route for %s:%s already exsits", strings.ToLower(requestType.String()), path))
		return
	}

	var pathParts = strings.Split(path, "/")
	var routeParts = make([]routePartInfo, 0)

	for _, part := range pathParts {
		var isVariable = strings.HasPrefix(part, ":")

		if isVariable {
			part = part[1:]
		}

		routeParts = append(routeParts, routePartInfo{
			isVariable, part,
		})
	}

	log.Verbose(_logTag, "route parts: %v", routeParts)

	_routes[path] = &routeInfo{
		requestType: requestType,
		isPublic:    isPublic,
		handler:     handler,
		parts:       routeParts,
		contract:    contract,
	}

	log.Information(_logTag, "added route for %s:%s [public: %v]", strings.ToLower(requestType.String()), path, isPublic)
}
