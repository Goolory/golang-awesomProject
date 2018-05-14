package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
	"strconv"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
)

func getAdminIdFromToken(c *gin.Context, db *gorm.DB) uint32 {
	var token string
	h := c.Request.Header.Get("Authorization")
	if h == "" {
		token = c.Request.URL.Query().Get("access_token")
	} else {
		ss := strings.Split(h, " ")
		if len(ss) < 2 {
			return 0
		}
		if strings.TrimSpace(ss[0]) != "Bearer" {
			return 0
		}
		token = strings.TrimSpace(ss[1])
	}
	if token != "" {
		ss := strings.Split(token, "|")
		if len(ss) < 2 {
			return 0
		}
		tid, _ := strconv.Atoi(ss[0])
		var t dbmodel.Token
		if err := db.Where("id = ?", tid).First(&t).Error; err != nil {
			logger.Error(err)
			return 0
		}
		if t.AccessToken == token {
			return t.UserId
		}
	}
	return 0
}

func GetAdmintId(c *gin.Context, db *gorm.DB) uint32 {
	adminId := getAdminIdFromToken(c, db)
	//if adminId == 0 {
	//	adminId = getAdminIdFromCookie(c, db)
	//}
	return adminId
}
