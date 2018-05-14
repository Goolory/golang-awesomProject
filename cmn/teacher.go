package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func AdminTeacherAllHandler(c *gin.Context) {
	var teachers []*dbmodel.Teacher
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.Find(&teachers).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"teachers": teachers}})

}

func AdminTeacherUpdateHandler(c *gin.Context) {
	type param struct {
		Id          uint32 `json:"id"`
		TeacherNo   string `json:"teacher_no"`
		TeacherName string `json:"teacher_name"`
		Password    string `json:"password"`
		Type        uint32 `json:"type"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	logger.Debug(p.Id)
	var teacher dbmodel.Teacher
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.First(&teacher, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	teacher.TeacherName.String = p.TeacherName
	teacher.TeacherNo = p.TeacherNo
	teacher.Password = p.Password
	teacher.Type = p.Type

	if err := db.Save(&teacher).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}

func AdminTeacherInfoHandler(c *gin.Context) {
	type param struct {
		Id uint32 `form:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var teacher dbmodel.Teacher
	if err := db.First(&teacher, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"teacher": teacher}})
}

func deleteTeacher(db *gorm.DB, id uint32) error {
	if err := db.Where("id = ?", id).Delete(&dbmodel.Teacher{}).Error; err != nil {
		logger.Error(err)
		return err
	} else {
		return nil
	}
}

func AdminTeacherDelHandler(c *gin.Context) {
	type param struct {
		Ids []uint32 `json:"ids"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin() //数据库事务
	if len(p.Ids) > 0 {
		for i := 0; i < len(p.Ids); i++ {
			if err := deleteTeacher(tx, p.Ids[i]); err != nil {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if err := tx.Model(&dbmodel.Class{}).Where("teacher_id=?", p.Ids[i]).Update("teacher_id", 0).Error; err != nil {
				logger.Error(err)
				return
			}
		}
	}
	if err := tx.Commit().Error; err != nil { //提交数据库事务
		logger.Error(err)
		return
	}

	tx = nil
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func AdminTeacherListHandler(c *gin.Context) {
	type param struct {
		TeacherName string `form:"teacher_name"`
		Page        uint32 `form:"page"`
		PageSize    uint32 `form:"page_size"`
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
	if p.TeacherName != "" {
		sql = sql.Where("teacher_name = ?", p.TeacherName)
	}
	if err := sql.Table("teacher").Count(&count).Error; err != nil {
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
	var teachers []*dbmodel.Teacher
	if err := sql.Table("teacher").Order("created_at desc").Find(&teachers).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"teachers": teachers, "total": count}})
}

func AdminTeacherAddHandler(c *gin.Context) {
	type param struct {
		TeacherNo   string `json:"teacher_no"`
		TeacherName string `json:"teacher_name"`
		Password    string `json:"password"`
		Type        uint32 `json:"type"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}

	if p.TeacherNo == "" || p.TeacherName == "" || p.Password == "" || p.Type == 0 {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	logger.Debug(p.TeacherName)

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var teacher dbmodel.Teacher
	if err := db.Where("`teacher_no` = ?", p.TeacherNo).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var dbParam dbmodel.Teacher
			dbParam.TeacherNo = p.TeacherNo
			dbParam.TeacherName.Valid = true
			dbParam.TeacherName.String = p.TeacherName
			dbParam.Password = p.Password
			dbParam.Type = p.Type
			if errCreate := db.Create(&dbParam).Error; errCreate != nil {
				logger.Error(errCreate)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else if teacher.TeacherName.String != "" {
		constant.ErrMsg(c, constant.UserAlreadyExists)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}
