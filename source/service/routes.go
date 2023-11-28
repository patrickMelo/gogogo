package service

import (
	"fmt"
	"gogogo/data"
	"gogogo/data/contract"
	"gogogo/log"
	"gogogo/requests"
	"strings"
)

func AddPublicPull(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Pull, path, true, handler, contract)
}

func AddPublicPush(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Push, path, true, handler, contract)

}

func AddPublicUpdate(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Update, path, true, handler, contract)

}

func AddPublicDelete(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Delete, path, true, handler, contract)

}

func AddPrivatePull(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Pull, path, false, handler, contract)

}

func AddPrivatePush(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Push, path, false, handler, contract)
}

func AddPrivateUpdate(path string, handler requests.Handler, contract *contract.Contract) {
	addRoute(requests.Update, path, false, handler, contract)
}

func AddPrivateDelete(path string, handler requests.Handler, contract *contract.Contract) {
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
	contract    *contract.Contract
}

var (
	_routes = make([]*routeInfo, 0)
)

func breakPath(path string) []string {
	var pathParts = strings.Split(path, "/")

	if strings.HasPrefix(path, "/") {
		pathParts = pathParts[1:]
	}

	return pathParts
}

func findRoute(requestType requests.Type, path string) *routeInfo {
	var pathParts = breakPath(path)

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

func addRoute(requestType requests.Type, path string, isPublic bool, handler requests.Handler, contract *contract.Contract) {
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

	_routes = append(_routes, &routeInfo{
		requestType: requestType,
		isPublic:    isPublic,
		handler:     handler,
		parts:       routeParts,
		contract:    contract,
	})

	log.Information(_logTag, "added route for %s:%s [public: %v]", strings.ToLower(requestType.String()), path, isPublic)
}

func extractRouteData(route *routeInfo, path string) (routeData data.GenericMap) {
	var pathParts = strings.Split(path, "/")[1:]

	routeData = data.NewGenericMap()

	for index, part := range route.parts {
		if !part.isVariable {
			continue
		}

		routeData.Set(part.part, pathParts[index])
	}

	return
}
