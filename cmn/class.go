package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type classParam struct {
	Id          uint32    `json:"id"`
	ClassName   string    `json:"class_name"`
	TeacherName string    `json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func AdminClassAllHandler(c *gin.Context) {
	var class []*dbmodel.Class
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.Find(&class).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"class": class}})
}

func AdminClassUpdateHandler(c *gin.Context) {
	type param struct {
		Id        uint32 `json:"id"`
		ClassName string `json:"class_name"`
		TeacherId uint32 `json:"teacher_id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	var class dbmodel.Class
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.First(&class, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	class.ClassName.String = p.ClassName
	class.TeacherId = p.TeacherId

	if err := db.Save(&class).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}

func AdminClassInfoHandler(c *gin.Context) {
	type param struct {
		Id uint32 `form:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var class dbmodel.Class
	if err := db.First(&class, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"class": class}})

}

func AdminClassDelHandler(c *gin.Context) {
	type param struct {
		Ids []uint32 `json:"ids"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	tx := db.Begin()
	if len(p.Ids) > 0 {
		for i := 0; i < len(p.Ids); i++ {
			if err := tx.Where("id = ?", p.Ids[i]).Delete(&dbmodel.Class{}).Error; err != nil {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if err := tx.Model(&dbmodel.User{}).Where("class_id = ?", p.Ids[i]).Update("class_id", 0).Error; err != nil {
				logger.Error(err)
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error(err)
		return
	}
	tx = nil
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}

func AdminClassListHandler(c *gin.Context) {
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
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	count := 0
	sql := db.Where("1=1")

	if err := sql.Table("class").Count(&count).Error; err != nil {
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

	var clazz []classParam
	if err := sql.Table("`class` as c").Joins("left join `teacher` as t on c.teacher_id = t.id").Order("c.created_at desc").Select("c.id, c.class_name, t.teacher_name, c.created_at, c.updated_at").Scan(&clazz).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"class": clazz, "total": count}})
}

func AdminClassAddHandler(c *gin.Context) {
	type param struct {
		ClassName string `json:"class_name"`
		TeacherId uint32 `json:"teacher_id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invaild request param", err)
		return
	}

	if p.ClassName == "" || p.TeacherId == 0 {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var class dbmodel.Class
	if err := db.Where("class_name = ?", p.ClassName).First(&class).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var dbP dbmodel.Class
			dbP.ClassName.String = p.ClassName
			dbP.ClassName.Valid = true
			dbP.TeacherId = p.TeacherId
			if errC := db.Create(&dbP).Error; errC != nil {
				logger.Error(errC)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else if class.ClassName.String != "" {
		constant.ErrMsg(c, constant.UserAlreadyExists)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}
