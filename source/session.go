package main

import (
	"gogogo/data"
	"gogogo/requests"

	"github.com/google/uuid"
)

var sessionLoginContract = data.NewContract([]data.ContractField{
	data.ContractString("username").Length(1, 128).Required(),
	data.ContractString("password").Length(8, 128).Required(),
})

func GetSessionLoginContract() *data.Contract {
	return sessionLoginContract
}

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
