package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id        uint32     `gorm:"primary_key;auto_increment" json:"id"`
	ClassId   uint32     `json:"class_id"`
	StudentNo NullString `gorm:"size:64" json:"student_no"`
	Username  NullString `gorm:"size:64" json:"username"`
	Password  string     `gorm:"size:64" json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

func InitUser(db *gorm.DB) error {
	var err error

	if db.HasTable(&User{}) {
		err = db.AutoMigrate(&User{}).Error
	} else {
		err = db.CreateTable(&User{}).Error
	}

	return err
}

func DropTableUser(db *gorm.DB) {
	db.DropTable(&User{})

}
