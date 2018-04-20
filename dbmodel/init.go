package dbmodel

import "github.com/jinzhu/gorm"

func InitDbModel(db *gorm.DB) {
	var err error
	err = InitUser(db)
	if err != nil {
		println("user dbmodel err", err)
		return
	}

	err = initFile(db)
	if err != nil {
		println("File dbmodel err", err)
		return
	}
}

func RebuildDbModel(db *gorm.DB) {
	DropTableUser(db)
	DropTableFile(db)
}
