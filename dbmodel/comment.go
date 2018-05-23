package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Comment struct {
	Id        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ParentId  uint32    `json:"parent_id"`
	ThemeId   uint32    `json:"theme_id"`
	UserId    uint32    `json:"user_id"`
	Publisher string    `gorm:"size:64" json:"publisher"`
	Content   string    `gorm:"type: longtext" json:"content"`
	Disabled  bool      `gorm:"default: false" json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Comment) TableName() string {
	return "comment"
}

func InitComment(db *gorm.DB) error {
	var err error

	if db.HasTable(&Comment{}) {
		err = db.AutoMigrate(&Comment{}).Error
	} else {
		err = db.CreateTable(&Comment{}).Error
	}

	return err
}

func DropTableComment(db *gorm.DB) {
	db.DropTable(&Comment{})

}
