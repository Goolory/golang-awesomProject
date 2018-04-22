package dbmodel

import (
	"awesomeProject/tool/uuid"
	"github.com/jinzhu/gorm"
)

type Class struct {
	Id       []byte `gorm:"primary_key; auto_increment" json:"id"`
	ClassName NullString `gorm:"size:64" json:"class_name"`
	Password string `gorm:"size: 255" json:"password"`
}

func (Class) TableName() string {
	return "class"
}

func (f *Class) BeforeCreate(scorpe *gorm.Scope) error {
	if f.Id == nil {
		scorpe.SetColumn("Id", ([]byte)(uuid.NewFastUUID()))
	}
	return nil
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
