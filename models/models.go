package models //与数据库交互

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// 存储数据库连接
var db *gorm.DB

// 结构体Model，它包含了一些公共的字段，用于被其他模型继承
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"ceated_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	// var (
	// 	err                                               error
	// 	dbType, dbName, user, password, host, tablePrefix string
	// )

	// sec, err := setting.Cfg.GetSection("database")
	// if err != nil {
	// 	log.Fatal(2, "Fail to get section 'database':  %v", err)
	// }

	// dbType = sec.Key("TYPE").String()
	// dbName = sec.Key("NAME").String()
	// user = sec.Key("USER").String()
	// password = sec.Key("PASSWORD").String()
	// host = sec.Key("HOST").String()
	// tablePrefix = sec.Key("TABLE_PREFIX").String()

	// db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName))

	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Println(err)
	}

	//设置了一个回调函数，用于自定义表名的生成规则
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

}

func CloseDB() {
	defer db.Close()
}

// 通过回调函数在创建记录时更新时间戳
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// scope.SetColumn(...) 假设没有指定 update_column 的字段，
// 我们默认在更新回调设置 ModifiedOn 的值
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		//检查表中是否存在"DeletedOn"字段，通常用于软删除功能
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		//只有在作用域限制处于启用状态且存在"DeletedOn"字段的情况下，才会进行软删除操作
		if !scope.Search.Unscoped && hasDeletedOnField {
			//构造一个UPDATE语句，将指定表中的"DeletedOn"字段更新为当前时间戳值
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// 用于在字符串前面添加一个空格
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
