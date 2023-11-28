package main

import (
	"gogogo/data/contract"
	"gogogo/requests"
	"regexp"

	"github.com/google/uuid"
)

type Note struct {
	Id       string `key:"id"`
	Name     string `key:"name"`
	Contents string `key:"contents"`
}

var (
	notes map[string]*Note = make(map[string]*Note)
)

func NotesPull(request *requests.Request, response *requests.Response) (err error) {
	for _, note := range notes {
		response.Data.Set(note.Id, note.Name)
	}

	return
}

var IdRegex = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

var NotePullContract = contract.New(contract.String("id").Regex(IdRegex).Required())

func NotePull(request *requests.Request, response *requests.Response) (err error) {
	var foundNote, noteExists = notes[request.Data.GetString("id", "")]

	if !noteExists {
		response.Status = requests.ResourceNotFound
		return
	}

	response.Data.FromStruct(foundNote)
	return
}

var NotePushContract = contract.New(
	contract.String("name").Length(3, 32).Required(),
	contract.String("contents").Length(1, 100).Optional(),
)

func NotePush(request *requests.Request, response *requests.Response) (err error) {
	var newNote = &Note{
		Id:       uuid.NewString(),
		Name:     request.Data.GetString("name", ""),
		Contents: request.Data.GetString("contents", ""),
	}

	notes[newNote.Id] = newNote

	response.Status = requests.ResourceCreated
	response.Data.Set("id", newNote.Id)
	return
}

var NoteUpdateContract = contract.New(
	contract.String("id").Regex(IdRegex).Required(),
	contract.String("name").Length(3, 32).Required(),
	contract.String("contents").Length(1, 100).Optional(),
)

func NoteUpdate(request *requests.Request, response *requests.Response) (err error) {
	var foundNote, noteExists = notes[request.Data.GetString("id", "")]

	if !noteExists {
		response.Status = requests.ResourceNotFound
		return
	}

	foundNote.Name = request.Data.GetString("name", "")
	foundNote.Contents = request.Data.GetString("contents", "")

	response.Status = requests.OK
	return
}
