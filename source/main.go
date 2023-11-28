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

	service.Start("GoGoGo", "0.1", listeners.Http())

	config.Load(config.Args(), config.Environment("GGG"), config.File())

	service.AddPublicPull("notes", NotesPull, nil)
	service.AddPublicPush("notes", NotePush, NotePushContract)
	service.AddPublicPull("notes/:id", NotePull, NotePullContract)
	service.AddPublicUpdate("notes/:id", NoteUpdate, NoteUpdateContract)

	service.Run()
}
