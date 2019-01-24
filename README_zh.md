# blog
基于 golang + vue 构建一个简单的个人博客系统

## 核心功能

### 主页面

* 显示个人基本信息
* 显示博客标签 (Python/C++等)
* 分页显示博客列表(博客标题, 博客摘要, 博客标签, 博客时间)

### 博客页面

* 显示博客标题
* 显示博客正文(基于 Markdown)
* 能够跳转到上一篇和下一篇

### 博客发布工具

并不准备搞一个管理页面. 借助一个控制台工具完成博客发布. 

### 搜索功能

按关键字对博客进行搜索. (TODO)

### 图床

方便处理 Markdown 中包含的图片. (TODO)

### 统计访问量

TODO

## 数据存储设计

使用 MySQL 作为数据存储媒介. 

### 表结构设计

**数据库名**

```sql
create database blog;
```

**博客表**

```sql
create table blogs (
    name varchar(255),
    description text,
    content text,
    create_time varchar(100),
    modify_time varchar(100),
    tag varchar(255)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;
```

其中 desc 描述信息, modify_time 修改时间, tag 博客标签暂时不考虑, 字段预留. 

content 中存的是博客正文的 html. 

### golang 操作 MySQL

依赖的包

```go
"database/sql"
_ "github.com/go-sql-driver/mysql"
```

基本操作(建立连接)

```go
Db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/blog?charset=utf8")
```

基本操作(插入数据)

```go
stmt, err := db.Prepare(`insert into blogs values(?, ?, ?, ?, ?, ?)`)
defer stmt.Close()
_, err = stmt.Exec(info.Html(), info.Description,
   info.Content, info.CreateTime, info.ModifyTime, info.Tag)
```

基本操作(查询)

```go
// 注意这一堆 defer 
stmt, err := Db.Prepare(`select name, create_time from blogs`)
defer stmt.Close()
rows, err := stmt.Query()
defer rows.Close()
for rows.Next() {
	blog := common.BlogInfo{}
	rows.Scan(&blog.Name, &blog.CreateTime)
}
```



## 博客发布工具

### Markdown 转换为 HTML

golang 中提供了一个 blackfriday 库, 可以很方便的完成 md 到 html 之间的转换

并且还提供了一个 blackfriday-tool (一个现成的命令行工具)

此处直接基于 blackfriday-tool 来完成转换. 

使用方法很简单

```sh
# 安装 blackfriday
go get -u gopkg.in/russross/blackfriday.v2

# 安装 blackfriday-tool
go get github.com/russross/blackfriday-tool

# 使用 blackfriday-tool 进行转换, 可以使用 -css 选项指定样式.
blackfriday-tool -css=sspai.css input.md output.html
```

其中 Markdown 样式文件出自

https://sspai.com/post/43873

### 获取环境变量
需要通过 GOBIN 来找到 blackfriday-tool 路径

```go
import "os"
gobin := os.Getenv("GOBIN")
```

### 创建子进程
通过子进程的方式调用 blackfriday-tool 
```go
import "os/exec"
bf := exec.Command(gobin+"/blackfriday-tool", "-css=sspai.css", inputPath, outputPath)
stdout, err := bf.Output()
```

## 博客服务器

### 博客服务API设计

使用 gin 在后端提供以下 Rest API

#### / 

获取第 1 页的博客列表. 

**响应示例**

```json
{
  "ok": true,
  "reason": "",
	"blogs": [
        {
            "name": "基于golang+vue的博客系统",
            "desc": "基于 golang + vue 构建一个简单的个人博客系统...",
            "datetime": "2019/01/14 14:07",
            "tag": ["golang", "vue"],
        },
     ],
}
```

#### /page/:page?month=[:month]&tag=[:tag]

获取第 :page 页的博客列表. 将所有的博客按照时间逆序排序. 每页10篇. 

支持过滤器. 获取某个月/某个标签下的博客

**响应示例**

同上

#### /blog/:id

获取指定 id 的博客详细信息

响应结果为一个 Markdown 转为的 html

**响应示例**

```json
{
    "ok": true,
    "reason": "",
    blog: {
        "Name": "基于golang+vue的博客系统.html",
        "Description": "",
        "Content": ".......",
        "CreateTime": "20190124",
        "ModifyTime": "20190124",
        "Tag": "golang",
    }
}
```

## 博客前端





