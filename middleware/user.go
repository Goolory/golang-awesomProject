package middleware

import (
	"github.com/gin-gonic/gin"
	"awesomeProject/constant"
	"github.com/jinzhu/gorm"
	"strings"
	"strconv"
	"awesomeProject/tool/logger"
	"awesomeProject/dbmodel"
	"net/http"
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
//func getAdminIdFromCookie(c *gin.Context, db *gorm.DB) uint32 {
//	v := cookie.GetVal(constant.CookieUserId, c.Request)
//	if v == nil {
//		return 0
//	}
//	switch v.(type) {
//	case string:
//		userId, _ := strconv.Atoi(v.(string))
//		if userId == 0 {
//			return 0
//		}
//		v = cookie.GetVal(constant.CookieExpire, c.Request)
//		switch v.(type) {
//		case string:
//			expire, err := time.Parse(time.RFC3339, v.(string))
//			if err != nil {
//				return 0
//			}
//			if expire.Unix() <= time.Now().Unix() {
//				return 0
//			}
//			if expire.Unix()-time.Now().Unix() < 15*60 {
//				cookie.SetVal(constant.CookieExpire, expire.Add(time.Second*10*60).Format(time.RFC3339), c.Writer)
//			}
//		}
//	}
//	return 0
//}


func getAdminId(c *gin.Context, db *gorm.DB) uint32 {
	adminId := getAdminIdFromToken(c, db)
	//if adminId == 0 {
	//	adminId = getAdminIdFromCookie(c, db)
	//}
	return adminId
}

func AdminVerifyHandler(c *gin.Context)  {
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	adminId := getAdminId(c, db)
	if adminId == 0 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	var admin dbmodel.Admin
	if err := db.Where("id = ?", adminId).First(&admin).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Set(constant.ContextAdmin, admin)
	c.Next()
}
