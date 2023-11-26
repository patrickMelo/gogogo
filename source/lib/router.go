package lib

import (
	"fmt"
	"strings"
)

func AddPublicPull(path string, handler RequestHandler, contract *Contract) {
	addRoute(PullRequest, path, true, handler, contract)
}

func AddPublicPush(path string, handler RequestHandler, contract *Contract) {
	addRoute(PushRequest, path, true, handler, contract)

}

func AddPublicUpdate(path string, handler RequestHandler, contract *Contract) {
	addRoute(UpdateRequest, path, true, handler, contract)

}

func AddPublicDelete(path string, handler RequestHandler, contract *Contract) {
	addRoute(DeleteRequest, path, true, handler, contract)

}

func AddPrivatePull(path string, handler RequestHandler, contract *Contract) {
	addRoute(PullRequest, path, false, handler, contract)

}

func AddPrivatePush(path string, handler RequestHandler, contract *Contract) {
	addRoute(PushRequest, path, false, handler, contract)
}

func AddPrivateUpdate(path string, handler RequestHandler, contract *Contract) {
	addRoute(UpdateRequest, path, false, handler, contract)
}

func AddPrivateDelete(path string, handler RequestHandler, contract *Contract) {
	addRoute(DeleteRequest, path, false, handler, contract)
}

type routePartInfo struct {
	isVariable bool
	part       string
}

type routeInfo struct {
	requestType RequestType
	isPublic    bool
	handler     RequestHandler
	parts       []routePartInfo
	contract    *Contract
}

const routerLogTag = "router"

var (
	rourterRoutes = make(map[string]*routeInfo)
)

func findRoute(requestType RequestType, path string) *routeInfo {
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

	for _, route := range rourterRoutes {
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

func addRoute(requestType RequestType, path string, isPublic bool, handler RequestHandler, contract *Contract) {
	LogVerbose(routerLogTag, "adding route for %s:%s", strings.ToLower(requestType.String()), path)

	if existingRoute := findRoute(requestType, path); existingRoute != nil {
		LogError(routerLogTag, fmt.Errorf("route for %s:%s already exsits", strings.ToLower(requestType.String()), path))
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

	LogVerbose(routerLogTag, "route parts: %v", routeParts)

	rourterRoutes[path] = &routeInfo{
		requestType: requestType,
		isPublic:    isPublic,
		handler:     handler,
		parts:       routeParts,
		contract:    contract,
	}

	LogInformation(routerLogTag, "added route for %s:%s [public: %v]", strings.ToLower(requestType.String()), path, isPublic)
}
