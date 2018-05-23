package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type homeworkList struct {
	dbmodel.Homework
	Title       string `json:"title"`
	StudentName string `json:"student_name"`
	ClassName   string `json:"class_name"`
}

func TeacherHomeworkListHandler(c *gin.Context) {
	type param struct {
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
	teacher := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	logger.Debug(teacher.Id)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	count := 0

	sql := db.Where("1=1")
	if err := sql.Table("homework as h").Joins("left join test as t on t.id = h.parent_id "+
		"left join student_class_teacher as u on u.id = h.user_id").Select(""+
		"h.*, t.title, u.username as student_name, u.class_name").Where("u.teacher_id = ?", teacher.Id).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Debug(count)

	if p.Page != 0 {
		sql = sql.Offset(int(p.Page-1) * int(p.PageSize))
	}
	if p.PageSize > 0 {
		sql = sql.Limit(int(p.PageSize))
	}
	var hList []homeworkList
	if err := sql.Table("homework as h").Joins("left join test as t on t.id = h.parent_id left join student_class_teacher as u on u.id = h.user_id").Select(""+
		"h.*, t.title, u.username as student_name, u.class_name").Where("u.teacher_id = ?", teacher.Id).Scan(&hList).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": hList, "count": count})
}

func StudentHomeworkAddHandler(c *gin.Context) {
	type param struct {
		ParentId uint32 `json:"parent_id"`
		FilePath string `json:"file_path"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invaild request param", err)
		return
	}

	if p.ParentId == 0 || p.FilePath == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	student := c.MustGet(constant.ContextStudent).(dbmodel.User)
	logger.Debug(student.Id)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var homework dbmodel.Homework
	homework.UserId = student.Id
	homework.ParentId = p.ParentId
	homework.FilePath = p.FilePath
	if errC := db.Create(&homework).Error; errC != nil {
		logger.Error(errC)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}
