package main

import (
	"fmt"
	"testing"
)

const Md = `---
createAt: 2022-04-24 09:06
updateAt: 2022-04-24 09:06:45
publish: false
---

tags: #daily #2022-04

<< [[2022-04-23]] | [[2022-04-25]] >>

## 重点关注
### ==早上 7 件事==
- [ ] 检查下[滴答清单](https://dida365.com/webapp/#q/all/today)
- [ ] 花点时间回顾和反思
- [ ] 查看「反向链接」和「工作待办」
- [ ] 扫一眼邮件
- [ ] 确定最困难的工作，拆分成多个小任务
- [ ] 写下需要思考的东西


	
## 阅读笔记 & 会议纪要
通常记录一些需要技术阅读的内容

## 间歇日记 事件
记录了一天的内容



## TODO
剩余未完成内容



![[Pasted image 20220424145346.png]]]

`

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
