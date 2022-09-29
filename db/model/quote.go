package model

import "time"

type Quote struct {
	Id        int
	Number    int
	Content   string
	Author    string
	Submitter string
	Channel   string
	Added     time.Time
}

type Message struct {
	Id       string
	Content  string
	Username string
	Channel  string
	Created  time.Time
}
