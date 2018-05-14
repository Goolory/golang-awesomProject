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

type studentParam struct {
	Id        uint32    `json:"id"`
	ClassName string    `json:"class_name"`
	StudentNo string    `json:"student_no"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func AdminStudentUpdateHandler(c *gin.Context) {
	type param struct {
		Id        uint32 `json:"id"`
		ClassId   uint32 `json:"class_id"`
		StudentNo string `json:"student_no"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	var student dbmodel.User
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.First(&student, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	student.ClassId = p.ClassId
	student.StudentNo.String = p.StudentNo
	student.Username.String = p.Username
	student.Password = p.Password

	if err := db.Save(&student).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func AdminStudentInfoHandler(c *gin.Context) {
	type param struct {
		Id uint32 `form:"id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var student dbmodel.User
	if err := db.First(&student, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"students": student}})
}

func AdminStudentDelHandler(c *gin.Context) {
	type param struct {
		Ids []uint32 `json:"ids"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	logger.Debug(p.Ids)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if len(p.Ids) > 0 {
		for i := 0; i < len(p.Ids); i++ {
			if err := db.Where("id = ?", p.Ids[i]).Delete(&dbmodel.User{}).Error; err != nil {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func AdminStudentListHandler(c *gin.Context) {
	type param struct {
		ClassId  uint32 `form:"class_id"`
		Username string `form:"username"`
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
	if p.Username != "" {
		sql = sql.Where("username = ?", p.Username)
	}
	if p.ClassId != 0 {
		sql = sql.Where("class_id = ?", p.ClassId)
	}
	if err := sql.Table("user").Count(&count).Error; err != nil {
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
	var students []studentParam
	if err := sql.Table("user as u").Joins("left join class as c on c.id = u.class_id").Order("created_at desc").Select("u.id, c.class_name, u.student_no, u.username, u.password, u.created_at, u.updated_at").Scan(&students).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"students": students, "total": count}})
}

func AdminStudentAddHandler(c *gin.Context) {
	type param struct {
		ClasssId  uint32 `json:"class_id"`
		StudentNo string `json:"student_no"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invaild request param", err)
		return
	}

	if p.ClasssId == 0 || p.StudentNo == "" || p.Username == "" || p.Password == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var student dbmodel.User
	if err := db.Where("student_no = ?", p.StudentNo).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var dbP dbmodel.User
			dbP.ClassId = p.ClasssId
			dbP.StudentNo.Valid = true
			dbP.StudentNo.String = p.StudentNo
			dbP.Username.String = p.Username
			dbP.Username.Valid = true
			dbP.Password = p.Password
			if errC := db.Create(&dbP).Error; errC != nil {
				logger.Error(errC)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else if student.StudentNo.String != "" {
		constant.ErrMsg(c, constant.UserAlreadyExists)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
