package main

import (
	"gogogo/data/contract"
	"gogogo/requests"

	"github.com/google/uuid"
)

var SessionLoginContract = contract.New(
	contract.String("username").Length(1, 128).Required(),
	contract.String("password").Length(8, 128).Required(),
)

type SessionLoginData struct {
	SessionId string `key:"sessionId"`
}

func SessionLogin(request *requests.Request, response *requests.Response) error {
	response.Data.FromStruct(SessionLoginData{
		SessionId: uuid.NewString(),
	})

	return nil
}

func SessionLogout(request *requests.Request, response *requests.Response) error {
	return nil
}
