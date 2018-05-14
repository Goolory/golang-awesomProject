package cmn

import (
	"github.com/gin-gonic/gin"
	"awesomeProject/tool/logger"
	"github.com/jinzhu/gorm"
	"awesomeProject/dbmodel"
	"net/http"
	"awesomeProject/constant"
	"awesomeProject/view"
)

func StudentQuestionListHandler(c *gin.Context) {
	type param struct {
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

	student := c.MustGet(constant.ContextStudent).(dbmodel.User)
	var sct view.StudentClassTeacher
	if err := db.First(&sct, student.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}


	sql := db.Where("1=1")

	sql = sql.Where("teacher_id = ?",sct.TeacherId)
	if err := sql.Table("question").Count(&count).Error; err != nil {
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

	var question []*dbmodel.Question
	if err := sql.Order("created_at desc").Find(&question).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"question": question, "total": count}})
}

func TeacherQuestionListHandler(c *gin.Context) {
	type param struct {
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
	teacher := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	sql := db.Where("1=1")

	sql = sql.Where("teacher_id = ?",teacher.Id)
	if err := sql.Table("question").Count(&count).Error; err != nil {
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

	var cList []*dbmodel.Question
	if err := sql.Order("created_at desc").Find(&cList).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"question": cList, "total": count}})
}

func TeacherQuestionAddHandler(c *gin.Context)  {
	type param struct {
		Title string `json:"title"`
		Content string `json:"content"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}

	if p.Content == "" || p.Title == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	teacher := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var question dbmodel.Question
	question.Title = p.Title
	question.Content = p.Content
	question.TeacherId = teacher.Id
	if errC := db.Create(&question).Error; errC != nil {
		logger.Error(errC)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
