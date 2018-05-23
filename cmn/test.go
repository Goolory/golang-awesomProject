package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"awesomeProject/view"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func TeacherTestDelHandler(c *gin.Context) {
	type param struct {
		Ids []uint32 `json:"ids"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if len(p.Ids) > 0 {
		for i := 0; i < len(p.Ids); i++ {
			if err := db.Where("id = ?", p.Ids[i]).Delete(&dbmodel.Test{}).Error; err != nil {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func TeacherTestUpdateHandler(c *gin.Context) {
	type param struct {
		Id uint32 `json:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	var test dbmodel.Test
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.First(&test, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	test.State = dbmodel.TestfinishedState

	if err := db.Save(&test).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func StudentTestInfoHandler(c *gin.Context) {
	type param struct {
		Id uint32 `form:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var test view.TestTeacher
	if err := db.First(&test, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"test": test}})
}
func TeacherTestInfoHandler(c *gin.Context) {
	type param struct {
		Id uint32 `form:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var test view.TestTeacher
	if err := db.First(&test, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"test": test}})
}

func TeacherTestListHandler(c *gin.Context) {
	type param struct {
		//UserId uint32 `form:"user_id"`
		//Username string `json:"username"`
		Page     uint32 `form:"page"`
		PageSize uint32 `form:"page_size"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	if p.Page <= 0 {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	count := 0

	sql := db.Where("1=1")
	teacher := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	//logger.Debug(p.UserId)
	if teacher.Id != 0 {
		sql = sql.Where("user_id = ?", teacher.Id)
	}
	if err := sql.Model(&view.TestTeacher{}).Count(&count).Error; err != nil {
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
	var tests []*view.TestTeacher
	if err := sql.Order("created_at desc").Find(&tests).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"tests": tests, "total": count}})
}
func TeacherTestAddHandler(c *gin.Context) {
	type param struct {
		UserId   uint32 `json:"user_id"`
		Title    string `json:"title"`
		Describe string `json:"discribe"`
		Content  string `json:"content"`
		FilesId  uint32 `json:"files_id"`
		View     string `json:"view"`
		Homework string `json:"homework"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invaild request param", err)
		return
	}
	logger.Debug(p.Title)
	logger.Debug(p.Describe)
	logger.Debug(p.Content)

	if p.UserId == 0 || p.Title == "" || p.Describe == "" || p.Content == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var test dbmodel.Test
	test.UserId = p.UserId
	test.Title = p.Title
	test.Describe = p.Describe
	test.Content = p.Content
	test.FilesId = p.FilesId
	test.View = p.View
	test.Homework = p.Homework
	if err := db.Create(&test).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "test_id": test.Id})
}

func StudentTestListHandler(c *gin.Context) {
	type param struct {
		//StudentId uint32 `form:"student_id"`
		Page     uint32 `form:"page"`
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
	logger.Debug(student.Id)

	var sct view.StudentClassTeacher
	if err := db.Where("id =?", student.Id).First(&sct).Error; err != nil {
		logger.Error(err)
		return
	}

	sql := db.Where("1=1")

	sql = sql.Where("user_id = ?", sct.TeacherId)

	if err := sql.Model(&view.TestTeacher{}).Count(&count).Error; err != nil {
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
	var tests []*view.TestTeacher
	if err := sql.Order("created_at desc").Find(&tests).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"tests": tests, "total": count}})
}
