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

func generaterTonkeS(student *dbmodel.User, db *gorm.DB) (*dbmodel.Token, error) {
	var token dbmodel.Token

	token.Disabled = false
	token.UserId = student.Id
	token.AccessToken = hashToken(student.StudentNo.String, student.Password, student.CreatedAt.String(), time.Now().String())
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

func StudentLoginHandler(c *gin.Context) {
	type param struct {
		StudentNo string `json:"student_no"`
		Password  string `json:"password"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	logger.Debug(p.StudentNo)

	if p.StudentNo == "" || p.Password == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var student dbmodel.User
	if erradmin := db.Where("`student_no` = ?", p.StudentNo).First(&student).Error; erradmin != nil {
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
	if student.Password != p.Password {
		constant.ErrMsg(c, constant.WrongPassword)
		return
	}

	token, err := generaterTonkeS(&student, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"access_token": token.AccessToken}})
}

func StudentInfoHandler(c *gin.Context) {
	student := c.MustGet(constant.ContextStudent).(dbmodel.User)
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"student_name": student.Username, "student_id": student.Id}})
}
