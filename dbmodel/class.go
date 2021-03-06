package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Class struct {
	Id        uint32     `gorm:"primary_key;auto_increment" json:"id"`
	ClassName NullString `gorm:"size:64" json:"class_name"`
	TeacherId uint32     `json:"teacher_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Class) TableName() string {
	return "class"
}

func InitClass(db *gorm.DB) error {
	var err error

	if db.HasTable(&Class{}) {
		err = db.AutoMigrate(&Class{}).Error
	} else {
		err = db.CreateTable(&Class{}).Error
	}

	return err
}

func DropTableClass(db *gorm.DB) {
	db.DropTable(&Class{})

}
