package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Teacher struct {
	Id          uint32     `gorm:"primary_key; auto_increment" json:"id"`
	TeacherNo   string     `gorm:"size 225" json:"teacher_no"`
	TeacherName NullString `gorm:"size:64" json:"teacher_name"`
	Password    string     `gorm:"size:64" json:"password"`
	Type        uint32     `json:"type"` //教授1 讲师2
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	Professor = (1 << 0)
	Lecturer  = (1 << 1)
)

func (Teacher) TableName() string {
	return "teacher"
}

func InitTeacher(db *gorm.DB) error {
	var err error

	if db.HasTable(&Teacher{}) {
		err = db.AutoMigrate(&Teacher{}).Error
	} else {
		err = db.CreateTable(&Teacher{}).Error
	}

	return err
}

func DropTableTeacher(db *gorm.DB) {
	db.DropTable(&Teacher{})

}
