package listeners

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gogogo/config"
	"gogogo/data"
	"gogogo/log"
	"gogogo/requests"
	"gogogo/service"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpListener struct {
	service.Listener

	err    error
	server http.Server
}

func Http() *HttpListener {
	return &HttpListener{}
}

func (listener *HttpListener) Start() (err error) {
	var keepAlive = config.GetBool("http.keepAlive", false)

	listener.server.Addr = config.GetString("http.listenAddress", ":80")
	listener.server.SetKeepAlivesEnabled(keepAlive)
	listener.server.Handler = &httpHandler{}

	log.Verbose(httpLogTag, "listen address = %s", listener.server.Addr)
	log.Verbose(httpLogTag, "keep alive = %v", keepAlive)

	log.Information(httpLogTag, "starting listener at '%s'", listener.server.Addr)

	go listener.asyncStart()
	time.Sleep(time.Millisecond * 500)

	return listener.err
}

func (listener *HttpListener) Stop() {
	log.Information(httpLogTag, "stopping")
	listener.server.Shutdown(context.Background())
}

func (listener *HttpListener) asyncStart() {
	listener.err = listener.server.ListenAndServe()

	if listener.err != http.ErrServerClosed {
		log.Error(httpLogTag, listener.err)
	}
}

const (
	httpLogTag = "http"
)

type httpHandler struct {
	http.Handler
}

func (handler *httpHandler) ServeHTTP(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	var request = requests.NewRequest()
	log.Information(httpLogTag, "(%s) %s %s", request.Id, httpRequest.Method, httpRequest.RequestURI)

	request.Type = requestTypeFromHttpMethod(httpRequest.Method)

	if request.Type == requests.Unknown {
		log.Warning(httpLogTag, "(%s) unsupported method: %s", request.Id, httpRequest.Method)
		httpResponse.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var writeError = func(status int, err error) {
		log.Error(httpLogTag, fmt.Errorf("(%s) %v", request.Id, err))

		httpResponse.WriteHeader(status)

		var jsonData, _ = json.Marshal(data.GenericMap{
			"error": fmt.Sprintf("%v", err),
		})

		httpResponse.Write(jsonData)
	}

	httpResponse.Header().Set("Content-Type", "application/json")

	log.Verbose(httpLogTag, "(%s) parsing headers", request.Id)

	if status, err := parseHeader(request, httpRequest.Header); err != nil {
		writeError(status, err)
		return
	}

	log.Verbose(httpLogTag, "(%s) parsing URL", request.Id)

	if err := parseUrl(request, httpRequest.URL); err != nil {
		writeError(http.StatusBadRequest, err)
		return
	}

	if (request.Type == requests.Push) || (request.Type == requests.Update) {
		log.Verbose(httpLogTag, "(%s) extracting body", request.Id)

		if err := extractBody(request, httpRequest.Body); err != nil {
			writeError(http.StatusUnprocessableEntity, err)
			return
		}
	}

	var response = requests.NewResponse(request.Id)
	var requestError = service.HandleRequest(request, response)

	if requestError != nil {
		writeError(http.StatusInternalServerError, requestError)
		return
	}

	log.Information(httpLogTag, "(%s) got %s with %d data entries", request.Id, response.Status, len(response.Data))

	httpResponse.WriteHeader(httpStatusFromResponseStatus(response.Status))
	var jsonData, _ = json.Marshal(response.Data)
	httpResponse.Write(jsonData)

	log.Verbose(httpLogTag, "(%s) %s", request.Id, base64.StdEncoding.EncodeToString(jsonData))

}

func parseHeader(request *requests.Request, header http.Header) (int, error) {
	var contentType = header.Get("Content-Type")

	log.Verbose(httpLogTag, "(%s) content type: %s", request.Id, contentType)

	if ((request.Type == requests.Push) || (request.Type == requests.Update)) && contentType != "application/json" {
		return http.StatusUnsupportedMediaType, fmt.Errorf("unsupported content-type: \"%s\"", contentType)
	}

	var tokenData = header.Get("Authorization")

	if tokenData != "" {
		log.Verbose(httpLogTag, "(%s) token data: %s", request.Id, tokenData)

		var splitData = strings.Split(tokenData, "Bearer ")

		if len(splitData) != 2 {
			return http.StatusBadRequest, fmt.Errorf("bad authorization token data: %s", tokenData)
		}

		request.Metadata.Set("token", splitData[1])
	}

	return http.StatusOK, nil
}

func parseUrl(request *requests.Request, url *url.URL) error {
	request.Path = url.EscapedPath()

	for key, values := range url.Query() {
		log.Verbose(httpLogTag, "(%s) %s = %v", request.Id, key, values)
		request.Data.Set(key, values)
	}

	return nil
}

func extractBody(request *requests.Request, body io.ReadCloser) (err error) {
	var bodyData []byte

	if bodyData, err = io.ReadAll(body); err != nil {
		return
	}

	log.Verbose(httpLogTag, "(%s) %s", request.Id, base64.StdEncoding.EncodeToString(bodyData))

	var decodedData = data.NewGenericMap()
	if err = json.Unmarshal(bodyData, &decodedData); err != nil {
		return
	}

	request.Data.MergeWith(decodedData)
	return
}

func requestTypeFromHttpMethod(httpMethod string) requests.Type {
	switch httpMethod {
	case http.MethodGet:
		return requests.Pull
	case http.MethodPost:
		return requests.Push
	case http.MethodPatch:
		return requests.Update
	case http.MethodDelete:
		return requests.Delete
	}

	return requests.Unknown
}

func httpStatusFromResponseStatus(status requests.Status) int {
	switch status {
	case requests.OK:
		return http.StatusOK
	case requests.InternalError:
		return http.StatusInternalServerError
	case requests.InvalidData:
		return http.StatusUnprocessableEntity
	case requests.NotAllowed:
		return http.StatusMethodNotAllowed
	case requests.NotAuthorized:
		return http.StatusForbidden
	case requests.AuthenticationRequired:
		return http.StatusUnauthorized
	case requests.ResourceCreated:
		return http.StatusCreated
	case requests.ResourceNotFound:
		return http.StatusNotFound
	case requests.ResourceAlreadyExists:
		return http.StatusFound
	}

	return http.StatusInternalServerError
}
