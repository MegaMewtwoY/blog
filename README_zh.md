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

### 博客内容

通过文件的方式存储博客内容. 

创建 ./data 目录, 作为数据的根目录. 

./data 目录中存放日期目录(精确到月, 例如 201901)

然后再在日期目录中存放博客(html文件).

对于每个博客的主键为 日期+html文件名字

目录结构如下:

```
# tree -N data/ 结构
data/
└── 基于golang+vue的博客系统.html
```

### 博客元数据

通过 Redis 的方式存储.

key 为 文件名(必须要求博客的名字不重复, 不过这样要求并不过分)

value 包含

* id (日期+html文件名字)
* 修改时间(YYYY/MM/DD HH:MM:SS)
* 标签(字符串列表)

## 博客发布工具

一个控制台工具, 能够将指定的 markdown 文件发布为博客

* 将 md 文件转为 html, 并发布到指定目录中.
* 修改 Redis 中的元数据

### Markdown 转换为 HTML

golang 中提供了一个 blackfriday 库, 可以很方便的完成 md 到 html 之间的转换

并且还提供了一个 blackfriday-tool (一个现成的命令行工具)

此处直接基于 blackfriday-tool 来完成转换. 

使用方法很简单

```
# 安装 blackfriday
go get -u gopkg.in/russross/blackfriday.v2

# 安装 blackfriday-tool
go get github.com/russross/blackfriday-tool

# 使用 blackfriday-tool 进行转换, 可以使用 -css 选项指定样式.
blackfriday-tool -css=sspai.css input.md output.html
```

其中 Markdown 样式文件出自

https://sspai.com/post/43873

### 操作 Redis 


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

略

## 博客前端




