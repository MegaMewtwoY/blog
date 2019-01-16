package main

import (
	"blog/common"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	outputPathPrefix = "../data/"
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
	// 3. 更新 Redis 中的数据
	err = updateRedis(blogInfo)
	if err != nil {
		fmt.Println("updateRedis failed!", err)
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
		ModifyDate: f.ModTime().Format("20060102"),
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
	return nil
}

func updateRedis(info *common.BlogInfo) error {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return fmt.Errorf("Redis connect failed! %s", err.Error())
	}
	defer c.Close()

	value, err := json.Marshal(*info)
	if err != nil {
		return fmt.Errorf("json.Marshal failed! %s", err.Error())
	}
	_, err = c.Do("SET", info.Name, value)
	if err != nil {
		return fmt.Errorf("Redis SET falied! %s", err.Error())
	}
	return nil
}
