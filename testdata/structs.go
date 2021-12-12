package testdata

import (
	"time"
)

type Author struct {
	ID       string `bson:"_id"`
	Name     string
	Password string `bson:"-"`
	Age      uint
	Books    []*Book
	Created  time.Time
	Deleted  *time.Time
}

type Book struct {
	ID    string `bson:"_id"`
	Name  string
	Pages uint
}
