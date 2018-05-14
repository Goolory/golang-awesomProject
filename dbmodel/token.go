package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Token struct {
	Id           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserId       uint32    `json:"user_id"`
	AccessToken  string    `gorm:"size:128" json:"access_token"`
	AccessExpire time.Time `json:"access_expire"`
	RefreshToken string    `gorm:"size:128" json:"refresh_token"`
	Disabled     bool      `json:"disabled"`
	//AccessAt     time.Time `json:"access_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Token) TableName() string {
	return "token"
}

func InitToken(db *gorm.DB) error {
	var err error

	if db.HasTable(&Token{}) {
		err = db.AutoMigrate(&Token{}).Error
	} else {
		err = db.CreateTable(&Token{}).Error
	}

	return err
}

func DropTableToken(db *gorm.DB) {
	db.DropTable(&Token{})

}
