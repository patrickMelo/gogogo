package lib

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (handler *httpHandler) ServeHTTP(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	var request = NewRequest()
	LogInformation(httpLogTag, "(%s) %s %s", request.Id, httpRequest.Method, httpRequest.RequestURI)

	request.Type = requestTypeFromHttpMethod(httpRequest.Method)

	if request.Type == UnknownRequest {
		LogWarning(httpLogTag, "(%s) unsupported method: %s", request.Id, httpRequest.Method)
		httpResponse.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var writeError = func(status int, err error) {
		LogError(httpLogTag, fmt.Errorf("(%s) %v", request.Id, err))

		httpResponse.WriteHeader(status)

		var jsonData, _ = json.Marshal(GenericMap{
			"error": fmt.Sprintf("%v", err),
		})

		httpResponse.Write(jsonData)
	}

	httpResponse.Header().Set("Content-Type", "application/json")

	LogVerbose(httpLogTag, "(%s) parsing headers", request.Id)

	if status, err := parseHeader(request, httpRequest.Header); err != nil {
		writeError(status, err)
		return
	}

	LogVerbose(httpLogTag, "(%s) parsing URL", request.Id)

	if err := parseUrl(request, httpRequest.URL); err != nil {
		writeError(http.StatusBadRequest, err)
		return
	}

	if (request.Type == PushRequest) || (request.Type == UpdateRequest) {
		LogVerbose(httpLogTag, "(%s) extracting body", request.Id)

		if err := extractBody(request, httpRequest.Body); err != nil {
			writeError(http.StatusUnprocessableEntity, err)
			return
		}
	}

	var response = NewResponse(request.Id)
	var requestError = handleRequest(request, response)

	if requestError != nil {
		writeError(http.StatusInternalServerError, requestError)
		return
	}

	LogInformation(httpLogTag, "(%s) got %s with %d data entries", request.Id, response.Status, len(response.Data))

	httpResponse.WriteHeader(httpStatusFromResponseStatus(response.Status))
	var jsonData, _ = json.Marshal(response.Data)
	httpResponse.Write(jsonData)

	LogVerbose(httpLogTag, "(%s) %s", request.Id, base64.StdEncoding.EncodeToString(jsonData))

}

type httpHandler struct {
	http.Handler
}

const httpLogTag = "http"

var (
	httpServer http.Server
	httpError  error
)

func startHttp(options GenericMap) bool {
	var keepAlive = options.GetBool("http.keepAlive", false)

	httpServer.Addr = options.GetString("http.listenAddress", ":80")
	httpServer.SetKeepAlivesEnabled(keepAlive)
	httpServer.Handler = &httpHandler{}

	LogVerbose(httpLogTag, "listen address = %s", httpServer.Addr)
	LogVerbose(httpLogTag, "keep alive = %v", keepAlive)

	LogInformation(httpLogTag, "starting listener at '%s'", httpServer.Addr)

	go asyncStartHttp()
	time.Sleep(time.Second)

	return httpError == nil
}

func stopHttp() {
	LogInformation(httpLogTag, "stopping")
	httpServer.Shutdown(context.Background())
}

func asyncStartHttp() {
	httpError = httpServer.ListenAndServe()

	if httpError != http.ErrServerClosed {
		LogError(httpLogTag, httpError)
	}
}

func parseHeader(request *Request, header http.Header) (int, error) {
	var contentType = header.Get("Content-Type")

	LogVerbose(httpLogTag, "(%s) content type: %s", request.Id, contentType)

	if ((request.Type == PushRequest) || (request.Type == UpdateRequest)) && contentType != "application/json" {
		return http.StatusUnsupportedMediaType, fmt.Errorf("unsupported content-type: \"%s\"", contentType)
	}

	var tokenData = header.Get("Authorization")

	if tokenData != "" {
		LogVerbose(httpLogTag, "(%s) token data: %s", request.Id, tokenData)

		var splitData = strings.Split(tokenData, "Bearer ")

		if len(splitData) != 2 {
			return http.StatusBadRequest, fmt.Errorf("bad authorization token data: %s", tokenData)
		}

		request.AuthToken = splitData[1]
	}

	return http.StatusOK, nil
}

func parseUrl(request *Request, url *url.URL) error {
	request.Path = url.EscapedPath()

	for key, values := range url.Query() {
		LogVerbose(httpLogTag, "(%s) %s = %v", request.Id, key, values)
		request.Data.Set(key, values)
	}

	return nil
}

func extractBody(request *Request, body io.ReadCloser) (err error) {
	var bodyData []byte

	if bodyData, err = io.ReadAll(body); err != nil {
		return
	}

	LogVerbose(httpLogTag, "(%s) %s", request.Id, base64.StdEncoding.EncodeToString(bodyData))

	var decodedData = NewGenericMap()
	if err = json.Unmarshal(bodyData, &decodedData); err != nil {
		return
	}

	request.Data.MergeWith(decodedData)
	return
}

func requestTypeFromHttpMethod(httpMethod string) RequestType {
	switch httpMethod {
	case http.MethodGet:
		return PullRequest
	case http.MethodPost:
		return PushRequest
	case http.MethodPatch:
		return UpdateRequest
	case http.MethodDelete:
		return DeleteRequest
	}

	return UnknownRequest
}

func httpStatusFromResponseStatus(status ResponseStatus) int {
	switch status {
	case OKStatus:
		return http.StatusOK
	case InternalErrorStatus:
		return http.StatusInternalServerError
	case InvalidDataStatus:
		return http.StatusUnprocessableEntity
	case NotAllowedStatus:
		return http.StatusMethodNotAllowed
	case NotAuthorizedStatus:
		return http.StatusForbidden
	case AuthenticationRequiredStatus:
		return http.StatusUnauthorized
	case ResourceCreatedStatus:
		return http.StatusCreated
	case ResourceNotFoundStatus:
		return http.StatusNotFound
	case ResourceAlreadyExistsStatus:
		return http.StatusFound
	}

	return http.StatusInternalServerError
}
