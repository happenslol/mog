package testdata

import (
	"time"

	"github.com/happenslol/mog/colgen"
)

type CustomID struct {
	Domain string
	ID     string
}

type PageCount int

type Author struct {
	ID       CustomID `bson:"_id"`
	Name     string
	Password string `bson:"-"`
	Age      uint
	Books    []*Book
	Created  time.Time
	Deleted  *time.Time
}

type Book struct {
	ID    CustomID `bson:"_id"`
	Name  string
	Pages PageCount
}
