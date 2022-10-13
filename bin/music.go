package bin

import (
	"io/fs"
)

type Music struct {
	AddedBy   string
	Title     string
	Channel   string
	Url       string
	Thumbnail string
	File      fs.File
}
