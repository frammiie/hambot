package model

import (
	"github.com/frammiie/hambot/db/types"
)

type Quote struct {
	Id        int
	Number    int
	Content   string
	Author    string
	Submitter string
	Channel   string
	Added     types.Timestamp
}

type Message struct {
	Id       string
	Content  string
	Username string
	Channel  string
	Created  types.Timestamp
}
