package requests

import (
	"crypto/sha512"
	"encoding/hex"
	"gogogo/data"
	"math/rand"
)

type Type int

const (
	Unknown Type = iota
	Pull
	Push
	Update
	Delete
)

func (_type Type) String() string {
	switch _type {
	case Unknown:
		return "Unknown"
	case Pull:
		return "Pull"
	case Push:
		return "Push"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	}

	return "?"
}

type Request struct {
	Id       string
	Type     Type
	Path     string
	Data     data.GenericMap
	Metadata data.GenericMap
}

// Creates a new, empty request (with a random ID).
func NewRequest() *Request {
	var randomBytes [64]byte

	for index := 0; index < 64; index++ {
		randomBytes[index] = byte(rand.Intn(255))
	}

	var randomHash = sha512.Sum512(randomBytes[:])
	var randomStart = rand.Intn(47)

	return &Request{
		Id:       hex.EncodeToString(randomHash[randomStart : randomStart+16]),
		Data:     data.NewGenericMap(),
		Metadata: data.NewGenericMap(),
	}
}

// The function signature for request handling functions.
type Handler func(request *Request, response *Response) error
