package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type File struct {
	Id        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	RelPath   string    `gorm:"size:255" json:"rel_path"`
	Filename  string    `gorm:"size:128" json:"filename"`
	Filesize  int64     `json:"filesize"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (File) TableName() string {
	return "file"
}

func initFile(db *gorm.DB) error {
	var err error

	if db.HasTable(&File{}) {
		err = db.AutoMigrate(&File{}).Error
	} else {
		err = db.CreateTable(&File{}).Error
	}

	return err
}

func DropTableFile(db *gorm.DB) {
	db.DropTableIfExists(&File{})
}
