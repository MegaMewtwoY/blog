package common

type BlogInfo struct {
	Name        string // 不包含扩展名 html
	Description string // 暂时不考虑
	Content     string // html 格式的博客正文
	CreateTime  string // 创建时间
	ModifyTime  string // 暂时不考虑
	Tag         string // 暂时不考虑
}

func (info BlogInfo) Markdown() string {
	return info.Name + ".md"
}

func (info BlogInfo) Html() string {
	return info.Name + ".html"
}

type Blogs []BlogInfo

func (blogs Blogs) Len() int {
	return len(blogs)
}

func (blogs Blogs) Less(i, j int) bool {
	return blogs[i].CreateTime > blogs[j].CreateTime
}

func (blogs Blogs) Swap(i, j int) {
	blogs[i], blogs[j] = blogs[j], blogs[i]
}
