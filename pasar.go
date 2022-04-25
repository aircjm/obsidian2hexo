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
const obsidianImgRegex = `!\[\[.*\]\]`

type HexoBlog struct {
	BlogList  *[]ObsidianNote
	ImageList []string
}

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
	o.FileName = filepath.Base(pathStr)
	o.FilePath = pathStr
	o.FileType = path.Ext(pathStr)
	o.Title = o.FileName[0 : len(o.FileName)-len(o.FileType)]
	return *o
}

func ConvertHexoBlog(list []ObsidianNote) {

	fileMap := make(map[string]ObsidianNote)

	for _, note := range list {
		fileMap[note.FileName] = note
	}

	for _, note := range list {
		if path.Ext(note.FilePath) != ".md" {
			continue
		}
		bytes, err := ioutil.ReadFile(note.FilePath)
		if err != nil {
			fmt.Errorf("处理失败")
		}
		// 开始解析yaml formatter
		md := string(bytes)
		m := regexp.MustCompile(hexoFormatterRegex).FindString(md)
		yml := FrontMatter{}
		err = yaml.Unmarshal([]byte(m), &yml)
		if err != nil {
			fmt.Println("转化yaml配置失败")
		}
		note.FrontMatter = yml
		// fmt.Println(note.FrontMatter)
		// 没有了formatter就是content
		content := regexp.MustCompile(hexoFormatterRegex).ReplaceAllString(md, "")
		note.Content = content

		// 查找图片附件的双向链接
		findAllString := regexp.MustCompile(obsidianImgRegex).FindAllString(md, -1)
		for _, s := range findAllString {
			fileFullName := s[3 : len(s)-2]
			fmt.Println(fileFullName)
			// 图片类型处理链接替换
			imgLink := "![" + fileFullName + "](./image/" + fileFullName + ")"
			note.MarkDown = strings.Replace(content, s, imgLink, -1)
		}
	}

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
