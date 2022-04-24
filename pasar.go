package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const hexoFormatterRegex = `(?i)(?m)(?s)---.*?---`

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

func readFile(pathStr string) ObsidianNote {
	o := new(ObsidianNote)
	o.Title = path.Base(pathStr)
	o.FileName = path.Base(pathStr)
	o.FilePath = pathStr
	o.FileType = path.Ext(pathStr)
	return *o
}

func GetMdList(list []ObsidianNote) []ObsidianNote {
	var mdList []ObsidianNote
	for _, note := range list {
		if path.Ext(note.FilePath) != ".md" {
			continue
		}

		bytes, err := ioutil.ReadFile(note.FilePath)
		if err != nil {
			fmt.Errorf("处理失败")
		}
		// 开始解析yaml formatter
		s := string(bytes)
		m := regexp.MustCompile(hexoFormatterRegex).FindString(s)
		yml := HexoFormatter{}
		err = yaml.Unmarshal([]byte(m), &yml)
		if err != nil {
			fmt.Println("转化yaml配置失败")
		}

	}
	return mdList
}

/**
---
createAt: 2022-04-20 20:04:05
updateAt: 2022-04-20 20:04:05
publish: false
draft: true
cards-deck: Default
---
*/
type HexoFormatter struct {
	CreateAt  string `yaml:"createAt"`
	UpdateAt  string `yaml:"updateAt"`
	Publish   bool   `yaml:"publish"`
	Draft     bool   `yaml:"draft"`
	CardsDeck string `yaml:"cards-deck"`
}
