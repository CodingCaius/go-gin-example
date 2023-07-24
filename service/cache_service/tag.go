package cache_service

import (
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

//用于获取标签列表在缓存中的键
//当不同的标签属性组合出现时，将生成不同的唯一缓存键，
//确保标签列表缓存的唯一性和准确性
func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}