package models

import (
	"github.com/jinzhu/gorm"

	"time"
)

type Article struct {
	Model

	//外键
	//gorm会通过类名+ID 的方式去找到这两个类之间的关联关系
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}


	return count, nil
}

// 分页获取符合指定条件的文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article

	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Related(&article.Tag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil

}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}
	return nil
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}



//硬删除要使用 Unscoped()，这是 GORM 的约定
func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}

