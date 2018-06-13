package models

import (
	"fmt"
)

//登录记录
type File struct {
	Model
	Name string `json:"name"`
	Path string `json:"path"`
	Url  string `json:"url"`
}

func GetFiles() []*File {
	files := []*File{}

	ret := db.Model(&File{}).Find(&files)
	if ret.Error != nil {
		fmt.Printf("[GetFiles] %s\n", ret.Error)
		return files
	}
	return files
}
