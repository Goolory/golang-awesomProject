package cmn

import (
	"github.com/gin-gonic/gin"
	"awesomeProject/tool/logger"
	"github.com/jinzhu/gorm"
	"awesomeProject/dbmodel"
	"net/http"
	"awesomeProject/constant"
)

func StudentAnswerAddHandler(c *gin.Context)  {
	type param struct {
		QuestionId uint32 `json:"question_id"`
		Content string `json:"content"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invaild request param", err)
		return
	}

	if p.QuestionId == 0 || p.Content == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	student := c.MustGet(constant.ContextStudent).(dbmodel.User)
	logger.Debug(student.Id)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var answer dbmodel.Answer
	answer.UserId = student.Id
	answer.QuestionId = p.QuestionId
	answer.Content = p.Content
	if errC := db.Create(&answer).Error; errC != nil {
		logger.Error(errC)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

type answerStudent struct{
	dbmodel.Answer
	Username string `json:"username"`
}

func TeacherAnswerListHandler(c *gin.Context)  {
	type param struct {
		QuestionId uint32 `form:"question_id"`
		Page uint32 `form:"page"`
		PageSize uint32 `form:"page_size"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	if p.Page <= 0 || p.PageSize <= 0 {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	count := 0
	sql := db.Where("1=1")

	sql = sql.Where("question_id = ?",p.QuestionId)
	if err := sql.Table("answer").Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if p.Page != 0 {
		sql = sql.Offset(int(p.Page-1) * int(p.PageSize))
	}
	if p.PageSize > 0 {
		sql = sql.Limit(int(p.PageSize))
	}

	var cList []answerStudent
	if err := sql.Table("answer as a").Joins("left join user as u on u.id = a.user_id").Select("a.*, u.username").Order("a.created_at desc").Scan(&cList).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"answer": cList, "total": count}})

}
