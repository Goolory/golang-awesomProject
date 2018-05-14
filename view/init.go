package view

import (
	"awesomeProject/tool/logger"
	"github.com/jinzhu/gorm"
)

type View interface {
	TableName() string
	CreateView(db *gorm.DB) error
}

func InitDbView(db *gorm.DB, views ...View) {
	for _, v := range views {
		if err := v.CreateView(db); err != nil {
			logger.Fatal("Init "+v.TableName()+" failed, ", err)
			return
		}
	}
}

func InitAllDbView(db *gorm.DB) {
	InitDbView(db, TestTeacher{},
		StudentClassTeacher{},
	)
}

func DropAllDbView(db *gorm.DB) {
	type view struct {
		TableName string
	}
	var views []*view

	// 获取当前Database的所有View
	if err := db.Table("information_schema.VIEWS").
		Select("TABLE_NAME as table_name").
		Where("TABLE_SCHEMA=?", db.Dialect().CurrentDatabase()).
		Find(&views).Error; err != nil {
		logger.Fatal("failed in getting views' name.")
	}

	for _, v := range views {
		db.Exec("DROP VIEW " + v.TableName)
	}
}
