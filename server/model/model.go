package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type BlogItem struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Date string `json:"date"`
	Path string `json:"path"`
}

type Blogs []BlogItem

func (blogs Blogs) Len() int {
	return len(blogs)
}

func (blogs Blogs) Less(i, j int) bool {
	return blogs[i].Date > blogs[j].Date
}

func (blogs Blogs) Swap(i, j int) {
	blogs[i], blogs[j] = blogs[j], blogs[i]
}

const (
	blogsPath = "../data/"
)

var (
	// 存储所有的博客数据
	blogs Blogs = make(Blogs, 0, 50)

	// 月份 -> 下标列表
	monthIndex map[string][]int = make(map[string][]int)

	// 标签 -> 下标列表
	tagIndex map[string][]int = make(map[string][]int)
)

func init() {
	load()
}

func load() {
	// 1. 打开 ../data 目录
	filepath.Walk(blogsPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		blogItem := BlogItem{
			Name: f.Name(),
			Path: path,
			Desc: getDesc(path),
			// 卧槽神奇了, 原来只能用 2006/1/2 15:04:05 进行格式化
			Date: f.ModTime().Format("2006/01"),
		}
		blogs = append(blogs, blogItem)
		return nil
	})
	sort.Sort(blogs)
	// 2. 构建两个倒排索引
	for i, blog := range blogs {
		if monthIndex[blog.Date] == nil {
			monthIndex[blog.Date] = make([]int, 0, 50)
		}
		monthIndex[blog.Date] = append(monthIndex[blog.Date], i)
	}
	// TODO tag 的设计再想想. 目录? 文件名?

	fmt.Println("Load", len(blogs), "items!")
}

func getDesc(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", path, "failed")
		return ""
	}
	defer f.Close()
	b := make([]byte, 50, 50)
	f.Read(b)
	// TODO 有可能导致只取到半个中文的情况.
	return string(b)
}

func GetPage(page string) []BlogItem {
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("page convert err! page =", page, err)
		return nil
	}
	// 判定页码是否合理
	if pageNum < 1 || pageNum > len(blogs)/10+1 {
		fmt.Println("page error! page =", page)
		return nil
	}
	// page 1 0-9
	// page 2 10-19
	// page 3 20-29
	// page n (n-1)*10 - n*10
	beg := (pageNum - 1) * 10 // 博客的起始下标
	end := pageNum * 10       // 博客的结束下标
	if end > len(blogs) {
		end = len(blogs)
	}
	return blogs[beg:end]
}
