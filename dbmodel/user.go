package dbmodel

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	Id       []byte `gorm:"primary_key;auto_increment" json:"id"`
	ClassId []byte	`json:"class_id"`
	Username NullString `gorm:"size:64" json:"username"`
	Password string `gorm:"size:64" json:"password"`
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
