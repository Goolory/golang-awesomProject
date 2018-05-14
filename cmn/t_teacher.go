package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

//生成token(令牌)
func generaterTonkeT(teacher *dbmodel.Teacher, db *gorm.DB) (*dbmodel.Token, error) {
	var token dbmodel.Token

	token.Disabled = false
	token.UserId = teacher.Id
	token.AccessToken = hashToken(teacher.TeacherNo, teacher.Password, teacher.CreatedAt.String(), time.Now().String())
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

func TeacherLoginHandler(c *gin.Context) {
	type param struct {
		TeacherNo string `json:"teacher_no"`
		Password  string `json:"password"`
	}

	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	logger.Debug(p.TeacherNo)

	if p.TeacherNo == "" || p.Password == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var teacher dbmodel.Teacher
	if erradmin := db.Where("`teacher_no` = ?", p.TeacherNo).First(&teacher).Error; erradmin != nil {
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
	if teacher.Password != p.Password {
		constant.ErrMsg(c, constant.WrongPassword)
		return
	}

	token, err := generaterTonkeT(&teacher, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"access_token": token.AccessToken}})
}

func TeacherInfoHandler(c *gin.Context) {
	teacher := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"teacher_name": teacher.TeacherName.String, "teacher_id": teacher.Id}})
}
