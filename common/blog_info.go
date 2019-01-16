package common

type BlogInfo struct {
	Name       string // 不包含扩展名
	ModifyDate string
	Desc       string   // 这个先不管
	Tag        []string // 这个也先不管
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
	return blogs[i].ModifyDate > blogs[j].ModifyDate
}

func (blogs Blogs) Swap(i, j int) {
	blogs[i], blogs[j] = blogs[j], blogs[i]
}
