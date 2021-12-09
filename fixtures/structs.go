package fixtures

import "time"

type Author struct {
	ID      string
	Name    string
	Age     uint
	Books   []*Book
	Created time.Time
	Deleted *time.Time
}

type Book struct {
	ID   string
	Name string
}
