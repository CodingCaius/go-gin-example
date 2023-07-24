//生成在缓存中存储文章信息所需的键（key）

package cache_service

import (
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
}

//用于获取单个文章在缓存中的键
func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

//用于获取文章列表在缓存中的键
func (a *Article) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}