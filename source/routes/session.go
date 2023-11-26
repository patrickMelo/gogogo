package routes

import (
	"gogogo/lib"

	"github.com/google/uuid"
)

var sessionLoginContract = lib.NewContract([]lib.ContractField{
	lib.ContractString("username").Length(1, 128).Required(),
	lib.ContractString("password").Length(8, 128).Required(),
})

func GetSessionLoginContract() *lib.Contract {
	return sessionLoginContract
}

type SessionLoginData struct {
	SessionId string `key:"sessionId"`
}

func SessionLogin(request *lib.Request, response *lib.Response) error {
	response.Data.FromStruct(SessionLoginData{
		SessionId: uuid.NewString(),
	})

	return nil
}

func SessionLogout(request *lib.Request, response *lib.Response) error {
	return nil
}
