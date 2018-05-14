package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"time"
)

func hashToken(ss ...string) string {
	s := strings.Join(ss, ":")
	h := sha512.Sum512([]byte(s))
	return base64.StdEncoding.EncodeToString(h[:])
}

//生成token(令牌)
func generaterTonke(admin *dbmodel.Admin, db *gorm.DB) (*dbmodel.Token, error) {
	var token dbmodel.Token

	token.Disabled = false
	token.UserId = admin.Id
	token.AccessToken = hashToken(admin.Username.String, admin.Password, admin.CreatedAt.String(), time.Now().String())
	token.AccessExpire = time.Now().AddDate(1, 0, 0)
	token.RefreshToken = hashToken(token.AccessToken, time.Now().String())
	if err := db.Create(&token).Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	token.AccessToken = fmt.Sprint(token.Id) + "|" + token.AccessToken

	if err := db.Model(&token).Update("access_token", token.AccessToken).Error; err != nil {
		logger.Error(err)
		db.Delete(&token)
		return nil, err
	}
	return &token, nil
}
func AdminRegisterHandler(c *gin.Context) {
	type param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}

	if p.Username == "" || p.Password == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var admin dbmodel.Admin
	if err := db.Where("`username` = ?", p.Username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var dbParam dbmodel.Admin
			dbParam.Username.Valid = true
			dbParam.Username.String = p.Username
			dbParam.Password = p.Password
			if createErr := db.Create(&dbParam).Error; createErr != nil {
				logger.Error(createErr)
				return
			}
		} else if err == nil && admin.Username.String != "" {
			constant.ErrMsg(c, constant.UserAlreadyExists)
		} else {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	token, err := generaterTonke(&admin, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"access_token": token.AccessToken}})

}
func AdminLoginHandler(c *gin.Context) {
	type param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}

	if p.Username == "" || p.Password == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var admin dbmodel.Admin
	if erradmin := db.Where("`username` = ?", p.Username).First(&admin).Error; erradmin != nil {
		if erradmin == gorm.ErrRecordNotFound {
			logger.Error(erradmin)
			constant.ErrMsg(c, constant.UserNotFound)
			return
		} else {
			logger.Error(erradmin)
			constant.ErrMsg(c, constant.StatusInternalServerError)
			return
		}
	}
	if admin.Password != p.Password {
		constant.ErrMsg(c, constant.WrongPassword)
		return
	}

	token, err := generaterTonke(&admin, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"access_token": token.AccessToken}})
}

func AdminInfoHandler(c *gin.Context) {
	admin := c.MustGet(constant.ContextAdmin).(dbmodel.Admin)
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"admin_name": admin.Username}})
}
