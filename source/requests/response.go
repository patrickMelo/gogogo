package requests

import "gogogo/data"

type Status int

const (
	OK Status = iota
	InternalError
	InvalidData
	NotAllowed
	NotAuthorized
	AuthenticationRequired
	ResourceCreated
	ResourceNotFound
	ResourceAlreadyExists
)

func (status Status) String() string {
	switch status {
	case OK:
		return "OK"
	case InternalError:
		return "InternalError"
	case InvalidData:
		return "InvalidData"
	case NotAllowed:
		return "NotAllowed"
	case NotAuthorized:
		return "NotAuthorized"
	case AuthenticationRequired:
		return "AuthenticationRequired"
	case ResourceCreated:
		return "ResourceCreated"
	case ResourceNotFound:
		return "ResourceNotFound"
	case ResourceAlreadyExists:
		return "ResourceAlreadyExists"
	}

	return "?"
}

type Response struct {
	RequestId string
	Status    Status
	Data      data.GenericMap
	Metadata  data.GenericMap
}

// Creates a new, empty response (using the specified request ID).
func NewResponse(requestId string) *Response {
	return &Response{
		RequestId: requestId,
		Data:      data.NewGenericMap(),
		Metadata:  data.NewGenericMap(),
		Status:    OK,
	}
}
