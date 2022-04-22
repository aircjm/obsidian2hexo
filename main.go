package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v2"
)


var AllNote  []ObsidianNote



func main()  {
	fmt.Println("开始处理")

	// 读取obsidian文件夹
	directory, err := LoadObsidianDirectory("C:\\Users\\admin\\Desktop\\lucida", nil, true)

	if err != nil {
		fmt.Errorf("读取文件出现错误")
	}
	println(len(AllNote))
	println(len(directory.Notes))
}


// FrontMatter is meta information for markdown documents
type FrontMatter map[string]interface{}

// ObsidianDirectory is a directory within an Obsidian Vault
type ObsidianDirectory struct {
	Name   string
	Path   string
	Parent *ObsidianDirectory
	Childs []ObsidianDirectory
	Notes  []ObsidianNote
	Files  []string
}


// ObsidianNote is a single note in Obsidian
type ObsidianNote struct {
	FrontMatter
	Title     string
	Content   string
	Directory *ObsidianDirectory
	FileName  string
	FilePath string
	FileType string
}



// LoadObsidianNote loads an Obsidian note from disk at given path
func LoadObsidianNote(path string) (ObsidianNote, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return ObsidianNote{}, err
	}

	matter, content, err := ParseFrontMatterMarkdown(raw)
	if err != nil {
		return ObsidianNote{}, err
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return ObsidianNote{
		FrontMatter: matter,
		Title:       title,
		Content:     content,
	}, nil
}


func ParseFrontMatterMarkdown(content []byte) (FrontMatter, string, error) {
	metaLines := make([]string, 0)
	bodyLines := make([]string, 0)
	state := 0

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if state < 2 && line == "---" {
			state++
			continue
		}
		if state == 1 {
			metaLines = append(metaLines, line)
		} else if state == 2 {
			bodyLines = append(bodyLines, line)
		}
	}
	if len(metaLines) == 0 {
		return nil, "", ErrNoFrontMatter
	}

	meta := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(strings.Join(metaLines, "\n")), &meta)
	if err != nil {
		return nil, "", err
	}

	return FrontMatter(meta), strings.TrimSpace(strings.Join(bodyLines, "\n")), nil
}
var (
	ErrNoFrontMatter = errors.New("missing front matter")
)

// ObsidianFilter includes or excludes a note
type ObsidianFilter func(ObsidianNote) bool



// LoadObsidianDirectory reads all notes and sub-directories within a directory in an Obsidian vault
func LoadObsidianDirectory(path string, filter ObsidianFilter, recurse bool) (root ObsidianDirectory, err error) {

	// 开始读文件夹
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	root.Path = path
	root.Name = filepath.Base(path)
	root.Childs = make([]ObsidianDirectory, 0)
	root.Files = make([]string, 0)
	root.Notes = make([]ObsidianNote, 0)


	for _, fi := range fis {

		if fi.IsDir() {
			fmt.Println("开始处理文件夹 ", fi.Name())
		}
		// ignore hidden
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		p := filepath.Join(path, fi.Name())

		// recurse directories
		if fi.IsDir() {
			if !recurse {
				continue
			}
			//log.WithField("directory", p).Debug("traverse sub-directory")
			sub, err := LoadObsidianDirectory(p, filter, true)
			if err != nil {
				return ObsidianDirectory{}, err
			} else if len(sub.Notes) == 0 {
				continue
			}

			sub.Parent = &root
			root.Childs = append(root.Childs, sub)

			// handle markdown files
		} else if filepath.Ext(p) == ".md" {
			//log.WithField("file", p).Debug("load markdown file")

			note, err := LoadObsidianNote(p)

			AllNote = append(AllNote, note)
			if err != nil {

				// ignore markdown files that lack front-matter
				if errors.Is(err, ErrNoFrontMatter) {
					//log.WithFields(log.Fields{"file": p}).Warn("ignore file with missing front matter")
					continue
				}
				return ObsidianDirectory{}, err
			}

			if filter != nil && !filter(note) {
				//log.WithField("note", note.Title).Info("note filtered out")
				continue
			}

			note.Directory = &root
			root.Notes = append(root.Notes, note)

			// handle other (static) files
		} else {
			//log.WithField("file", p).Debug("add static file")
			root.Files = append(root.Files, fi.Name())

		}
	}

	return
}


