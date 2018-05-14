package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Homework struct {
	Id        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ParentId  uint32    `json:"parent_id"`
	UserId    uint32    `json:"user_id"`
	FilePath string `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Homework) TableName() string {
	return "homework"
}

func initHomework(db *gorm.DB) error {
	var err error

	if db.HasTable(&Homework{}) {
		err = db.AutoMigrate(&Homework{}).Error
	} else {
		err = db.CreateTable(&Homework{}).Error
	}

	return err
}

func DropTableHomework(db *gorm.DB) {
	db.DropTableIfExists(&Homework{})
}
