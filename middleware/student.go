package middleware

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func getStudentIdFromToken(c *gin.Context, db *gorm.DB) uint32 {
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

func getStudentId(c *gin.Context, db *gorm.DB) uint32 {
	adminId := getStudentIdFromToken(c, db)
	//if adminId == 0 {
	//	adminId = getAdminIdFromCookie(c, db)
	//}
	return adminId
}

func StudentVerifyHandle(c *gin.Context) {
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	adminId := getStudentId(c, db)
	if adminId == 0 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	var student dbmodel.User
	if err := db.Where("id = ?", adminId).First(&student).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Set(constant.ContextStudent, student)
	c.Next()
}
