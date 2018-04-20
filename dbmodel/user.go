package dbmodel

import (
	"awesomeProject/tool/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id       []byte `gorm:"id"`
	Username string `gorm:"username"`
	Password string `gorm:"password"`
}

func (User) TableName() string {
	return "user"
}

func (f *User) BeforeCreate(scorpe *gorm.Scope) error {
	if f.Id == nil {
		scorpe.SetColumn("Id", ([]byte)(uuid.NewFastUUID()))
	}
	return nil
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
