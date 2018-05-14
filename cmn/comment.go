package cmn

import (
	"github.com/gin-gonic/gin"
	"awesomeProject/tool/logger"
	"awesomeProject/constant"
	"github.com/jinzhu/gorm"
	"net/http"
	"awesomeProject/dbmodel"
)


func TeacherCommentListHandler(c *gin.Context)  {
	StudentCommentListHandler(c)
}

func TeacherCommentAddHandler(c *gin.Context)  {
	user := c.MustGet(constant.ContextTeacher).(dbmodel.Teacher)
	addComment(c, user.Id, user.TeacherName.String)
}

func StudentCommentAddHandler(c *gin.Context)  {

	user := c.MustGet(constant.ContextStudent).(dbmodel.User)
	addComment(c, user.Id, user.Username.String)

}

func addComment(c *gin.Context, userId uint32, publisher string)  {
	type param struct {
		ThemeId uint32 `json:"theme_id"`
		Content string `json:"content"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}

	if p.Content == "" {
		constant.ErrMsg(c, constant.BadParameter)
		return
	}
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var comment dbmodel.Comment
	comment.ThemeId = p.ThemeId
	comment.UserId = userId
	comment.Publisher = publisher
	comment.Content = p.Content
	if errC := db.Create(&comment).Error; errC != nil {
		logger.Error(errC)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})

}

func AdminCommentDelHandler(c *gin.Context)  {
	type param struct {
		Id uint32 `json:"id"`
	}
	var p param
	if err := c.BindJSON(&p); err != nil {
		logger.Info("Invalid request param ", err)
		return
	}
	logger.Debug(p.Id)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if p.Id != 0 {
		if err := db.Where("id = ?", p.Id).Delete(&dbmodel.Comment{}).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
func AdminCommentListHandler(c *gin.Context)  {
	StudentCommentListHandler(c)
}

func StudentCommentListHandler(c *gin.Context) {
	type param struct {
		ThemeId uint32 `form:"theme_id"`
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

	sql = sql.Where("theme_id = ?",p.ThemeId)
	if err := sql.Table("comment").Count(&count).Error; err != nil {
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

	var cList []*dbmodel.Comment
	if err := sql.Order("created_at desc").Find(&cList).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": map[string]interface{}{"comments": cList, "total": count}})


}
