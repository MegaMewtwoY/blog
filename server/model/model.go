package model

import (
	"blog/common"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sort"
	"strconv"
)

var (
	Db *sql.DB
)

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/blog?charset=utf8")
	if err != nil {
		panic(fmt.Errorf("sql.Open failed! %s", err.Error()))
	}
	// 何时 Db.Close() ? 这是个问题
}

func GetPage(page string) (common.Blogs, error) {
	// 1. 获取到 http 请求中的页码
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("page convert err! page = %d, %s",
			pageNum, err.Error())
	}
	// 2. 预编译 sql
	stmt, err := Db.Prepare(`select name, create_time from blogs`)
	if err != nil {
		return nil, fmt.Errorf("Db.Prepare failed! %s", err.Error())
	}
	defer stmt.Close()
	// 3. 执行查询
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("stmt.Query failed! %s", err.Error())
	}
	defer rows.Close() // 这一堆 defer 是让我充满了 深深的恶意~~
	// 4. 获取结果
	blogs := make(common.Blogs, 0, 10)
	for rows.Next() {
		blog := common.BlogInfo{}
		rows.Scan(&blog.Name, &blog.CreateTime)
		blogs = append(blogs, blog)
	}
	// 5. 校验页码是否在有效范围
	//    例如 pageNum 为 2, 则如果博客数目 <= 10 篇就是非法的
	if pageNum < 1 || (pageNum-1)*10 >= len(blogs) {
		return nil, fmt.Errorf("pageNum not vaild! pageNum=%d", pageNum)
	}
	// 6. 获取本页要返回的博客序列
	beg := (pageNum - 1) * 10
	end := pageNum * 10
	if end > len(blogs) {
		end = len(blogs)
	}
	sort.Sort(blogs)
	return blogs[beg:end], nil
}
