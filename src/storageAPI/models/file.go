package models

import (
	"strconv"
	"time"
)

var (
	Storage map[string]*FileList
)

func init() {
	Storage = make(map[string]*FileList)
}

type File struct {
	id       string // must be lower case
	Name     string
	OwnerId  string
	IsPublic bool
}

type FileList struct {
	UserId string
	Files  map[string]*File
}

func AddFile(f File) string {
	f.id = "file_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	userId := f.OwnerId
	var l *FileList
	var ok bool
	if l, ok = Storage[userId]; !ok {
		l = &FileList{
			UserId: userId,
			Files:  make(map[string]*File),
		}
	}
	l.Files[f.id] = &f
	Storage[userId] = l
	return f.id
}
