package model

import (
	"blog/common"
	"fmt"
	"strconv"
)

const (
	blogsPath = "../data/"
)

func init() {

}

func GetPage(page string) common.Blogs {
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("page convert err! page =", page, err)
		return nil
	}
	return nil
}
