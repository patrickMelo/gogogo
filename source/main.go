package main

import (
	"gogogo/lib"
	"gogogo/routes"
)

func main() {
	lib.EnableVerboseLog()

	lib.Start("GoGoGo", "0.1")

	lib.AddPublicPush("session/login", routes.SessionLogin, routes.GetSessionLoginContract())
	lib.AddPrivatePush("session/logout", routes.SessionLogout, nil)

	lib.Run()
}
