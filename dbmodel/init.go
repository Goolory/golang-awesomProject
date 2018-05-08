package dbmodel

import (
	"awesomeProject/tool/logger"
	"github.com/jinzhu/gorm"
)

func InitDbModel(db *gorm.DB) {
	var err error
	err = InitAdmin(db)
	if err != nil {
		logger.Error("admin dbmodel err", err)
		return
	}
	err = InitUser(db)
	if err != nil {
		logger.Error("user dbmodel err", err)
		return
	}

	err = initFile(db)
	if err != nil {
		logger.Error("File dbmodel err", err)
		return
	}

	err = InitClass(db)
	if err != nil {
		logger.Error("class dbmodel err", err)
		return
	}
	err = InitToken(db)
	if err != nil {
		logger.Error("token model err", err)
		return
	}
	err = InitTeacher(db)
	if err != nil {
		logger.Error(err)
		return
	}

}

func RebuildDbModel(db *gorm.DB) {
	DropTableAdmin(db)
	DropTableUser(db)
	DropTableFile(db)
	DropTableToken(db)
	DropTableTeacher(db)
}
