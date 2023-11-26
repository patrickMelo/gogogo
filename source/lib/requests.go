package lib

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

type RequestType int

const (
	UnknownRequest RequestType = iota
	PullRequest
	PushRequest
	UpdateRequest
	DeleteRequest
)

func (_type RequestType) String() string {
	switch _type {
	case UnknownRequest:
		return "Unknown"
	case PullRequest:
		return "Pull"
	case PushRequest:
		return "Push"
	case UpdateRequest:
		return "Update"
	case DeleteRequest:
		return "Delete"
	}

	return "?"
}

type Request struct {
	Id        string
	Type      RequestType
	Path      string
	AuthToken string
	Data      GenericMap
}

// Creates a new, empty request (with a new unique ID).
func NewRequest() *Request {
	var randomBytes [64]byte

	for index := 0; index < 64; index++ {
		randomBytes[index] = byte(rand.Intn(255))
	}

	var randomHash = sha512.Sum512(randomBytes[:])
	var randomStart = rand.Intn(47)

	return &Request{
		Id:   hex.EncodeToString(randomHash[randomStart : randomStart+16]),
		Data: NewGenericMap(),
	}
}

// The function signature for request handling functions.
type RequestHandler func(request *Request, response *Response) error

type ResponseStatus int

const (
	OKStatus ResponseStatus = iota
	InternalErrorStatus
	InvalidDataStatus
	NotAllowedStatus
	NotAuthorizedStatus
	AuthenticationRequiredStatus
	ResourceCreatedStatus
	ResourceNotFoundStatus
	ResourceAlreadyExistsStatus
)

func (status ResponseStatus) String() string {
	switch status {
	case OKStatus:
		return "OK"
	case InternalErrorStatus:
		return "InternalError"
	case InvalidDataStatus:
		return "InvalidData"
	case NotAllowedStatus:
		return "NotAllowed"
	case NotAuthorizedStatus:
		return "NotAuthorized"
	case AuthenticationRequiredStatus:
		return "AuthenticationRequired"
	case ResourceCreatedStatus:
		return "ResourceCreated"
	case ResourceNotFoundStatus:
		return "ResourceNotFound"
	case ResourceAlreadyExistsStatus:
		return "ResourceAlreadyExists"
	}

	return "?"
}

type Response struct {
	RequestId string
	Status    ResponseStatus
	Data      GenericMap
}

// Creates a new, empty response (using the specified request ID).
func NewResponse(requestId string) *Response {
	return &Response{
		RequestId: requestId,
		Data:      NewGenericMap(),
		Status:    OKStatus,
	}
}
