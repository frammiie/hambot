package db

import "time"

type Quote struct {
	Number    int
	Content   string
	Author    string
	Submitter string
	Added     time.Time
}
