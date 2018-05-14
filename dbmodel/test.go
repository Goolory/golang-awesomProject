package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Test struct {
	Id        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserId    uint32    `json:"user_id"`
	Title     string    `gorm:"size:255" json:"title"`
	Describe  string    `json:"describe"`
	Content   string    ` json:"content"`
	FilesId   uint32    `json:"files_id"`
	View      string    `gorm:"size:225" json:"view"`
	Homework  string    `gorm:"size:225" json:"homework"`
	HomeworkState uint32 `gorm:"default:1" json:"homework_state"`
	State     uint32    `gorm:"default:'0'" json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	TestUnfinishedState = 0 //实验状态 未完成
	TestfinishedState   = 1 //实验状态 已完成

	HomeworkStateUnfinishedState = (1 < 0)
	HomeworkStateFinishedState = (1 < 1)
)

func (Test) TableName() string {
	return "test"
}

func InitTest(db *gorm.DB) error {
	var err error

	if db.HasTable(&Test{}) {
		err = db.AutoMigrate(&Test{}).Error
	} else {
		err = db.CreateTable(&Test{}).Error
	}

	return err
}

func DropTableTest(db *gorm.DB) {
	db.DropTable(&Test{})

}
