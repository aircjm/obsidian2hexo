package main

import (
	"fmt"
	"testing"
)

func TestLoadDir(t *testing.T) {
	var list []ObsidianNote
	LoadDir("C:\\Users\\admin\\Desktop\\lucida", &list)

	fmt.Println(len(list))
}

func TestGetMdList(t *testing.T) {

	var list []ObsidianNote
	LoadDir("C:\\Users\\admin\\Desktop\\lucida", &list)
	GetMdList(list)

}
