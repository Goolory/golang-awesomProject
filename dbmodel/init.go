package dbmodel

import (
	"github.com/jinzhu/gorm"
	"awesomeProject/tool/logger"
)

func InitDbModel(db *gorm.DB) {
	var err error
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
}

func RebuildDbModel(db *gorm.DB) {
	DropTableUser(db)
	DropTableFile(db)
}
