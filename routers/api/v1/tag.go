package v1

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
	// name := c.Query("name")

	// maps := make(map[string]interface{}) //存储查询条件
	// data := make(map[string]interface{}) //存储响应数据

	// if name != "" {
	// 	maps["name"] = name
	// }

	// var state int = -1 //state参数表示标签的特定状态
	// if arg := c.Query("state"); arg != "" {
	// 	state = com.StrTo(arg).MustInt()
	// 	maps["state"] = state
	// }

	// code := e.SUCCESS

	// data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	// data["total"] = models.GetTagTotal(maps)

	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": data,
	// })

}

// 新增文章标签
func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章标签
func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.PostForm("modified_by")

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         id,
		Name:       name,
		ModifiedBy: modifiedBy,
		State:      state,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
