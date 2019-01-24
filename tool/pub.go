package main

import (
	"blog/common"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	outputPathPrefix = "."
)

func main() {
	inputPath := ""
	flag.StringVar(&inputPath, "input", "", "待发布的博客")
	flag.Parse() // 这个必须要调用, 否则获取不到命令行参数
	// 1. 获取到这个博客的各种属性
	blogInfo := getBlogInfo(inputPath)
	if blogInfo == nil {
		fmt.Println("无法打开博客文件!", inputPath)
		os.Exit(1)
	}
	// 2. 进行 markdown -> html 的转换
	err := convertHtml(blogInfo, inputPath)
	if err != nil {
		fmt.Println("convertHtml failed!", err)
		os.Exit(1)
	}
	// 3. 更新 MySQL 中的数据
	err = updateMySQL(blogInfo)
	if err != nil {
		fmt.Println("updateMySQL failed!", err)
		os.Exit(1)
	}
	return
}

func getBlogInfo(inputPath string) *common.BlogInfo {
	f, err := os.Stat(inputPath)
	if err != nil {
		return nil
	}
	ext := path.Ext(f.Name())                     // 获取文件扩展名 .md
	fileName := strings.TrimSuffix(f.Name(), ext) // 去除扩展名, 只保留文件名
	blogInfo := common.BlogInfo{
		Name:       fileName,
		CreateTime: f.ModTime().Format("20060102"),
	}
	return &blogInfo
}

func convertHtml(info *common.BlogInfo, inputPath string) error {
	// 1. 构建输出文件名
	outputPath := outputPathPrefix + info.Html()
	// 2. 调用 blackfriday-tool 来完成转换
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		return fmt.Errorf("Not set env GOBIN!")
	}
	bf := exec.Command(gobin+"/blackfriday-tool", "-css=sspai.css", inputPath, outputPath)
	_, err := bf.Output()
	if err != nil {
		return fmt.Errorf("convertHtml failed! %s\n", err.Error())
	}
	// 3. 将刚才吐出的文件加载到 info 对象中
	html, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("ReadFile failed! %s", err.Error())
	}
	info.Content = string(html)
	return nil
}

func updateMySQL(info *common.BlogInfo) error {
	// 1. 连接数据库
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/blog?charset=utf8")
	if err != nil {
		return fmt.Errorf("MySQL Open failed")
	}
	defer db.Close()

	// 2. 插入数据
	stmt, err := db.Prepare(`insert into blogs values(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("db.Prepare failed!")
	}
	defer stmt.Close()
	_, err = stmt.Exec(info.Html(), info.Description,
		info.Content, info.CreateTime, info.ModifyTime, info.Tag)
	if err != nil {
		return fmt.Errorf("stmt.Exec failed")
	}
	return nil
}
