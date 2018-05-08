package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Admin struct {
	Id        uint32     `gorm:"primary_key;auto_increment" json:"id"`
	Username  NullString `gorm:"size:64" json:"username"`
	Password  string     `gorm:"size:64" json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admin"
}

func InitAdmin(db *gorm.DB) error {
	var err error

	if db.HasTable(&Admin{}) {
		err = db.AutoMigrate(&Admin{}).Error
	} else {
		err = db.CreateTable(&Admin{}).Error
	}

	return err
}

func DropTableAdmin(db *gorm.DB) {
	db.DropTable(&Admin{})

}
