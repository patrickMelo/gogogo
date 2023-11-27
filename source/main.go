package main

import (
	"gogogo/config"
	"gogogo/log"
	"gogogo/service"
	"gogogo/service/listeners"
)

func main() {
	log.Initialize(log.Stdout())
	log.SetMaxLevel(log.VerboseLevel)

	service.Start("GoGoGo", "0.1", []service.Listener{listeners.Http()})

	config.Load([]config.Provider{config.Args(), config.Environment("GGG"), config.File()})

	service.AddPublicPush("session/login", SessionLogin, GetSessionLoginContract())
	service.AddPrivatePush("session/logout", SessionLogout, nil)

	service.Run()
}
