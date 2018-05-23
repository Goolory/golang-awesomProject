package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Question struct {
	Id        uint32    `gorm:"primary_key; auto_increment" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TeacherId uint32    `json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Question) TableName() string {
	return "question"
}

func InitQuestion(db *gorm.DB) error {
	var err error

	if db.HasTable(&Question{}) {
		err = db.AutoMigrate(&Question{}).Error
	} else {
		err = db.CreateTable(&Question{}).Error
	}

	return err
}

func DropTableQuestion(db *gorm.DB) {
	db.DropTable(&Question{})

}
