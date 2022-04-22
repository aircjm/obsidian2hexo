package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)




func LoadDir(path string, list *[]ObsidianNote) {
	// 开始读文件夹
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, fileInfo := range fis {
		// ignore hidden
		if strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		p := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			LoadDir(p, list)
		} else {
			// 解析对应的markdown
			o := readFile(p)
			*list = append(*list, o)
		}
	}
}

func readFile(pathStr string)  ObsidianNote {
	o := new(ObsidianNote)
	o.Title = path.Base(pathStr)
	o.FileName = path.Base(pathStr)
	o.FilePath = pathStr
	o.FileType = path.Ext(pathStr)
	return *o
}



func GetMdList(list []ObsidianNote) []ObsidianNote {
	var mdList  []ObsidianNote

	for _, note := range list {
		if path.Ext(note.FilePath) != ".md" {
			continue
		}

		bytes, err := ioutil.ReadFile(note.FilePath)
		if err != nil {
			fmt.Errorf("处理失败")
		}
		s := string(bytes)
		m := regexp.MustCompile(`(?i)(?m)(?s)---.*?---`).ReplaceAllString(s, "")
		fmt.Println(m)
	}


	return mdList
}

