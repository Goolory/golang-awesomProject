package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Answer struct {
	Id          uint32     `gorm:"primary_key; auto_increment" json:"id"`
	QuestionId uint32 `json:"question_id"`
	UserId uint32 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}


func (Answer) TableName() string {
	return "answer"
}

func InitAnswer(db *gorm.DB) error {
	var err error

	if db.HasTable(&Answer{}) {
		err = db.AutoMigrate(&Answer{}).Error
	} else {
		err = db.CreateTable(&Answer{}).Error
	}

	return err
}

func DropTableAnswer(db *gorm.DB) {
	db.DropTable(&Answer{})

}
