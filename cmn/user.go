package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func UserRegisterHandler(c *gin.Context) {
	type param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var p param
	var err error

	if err = c.BindJSON(&p); err != nil {
		println("cmn/user.go invalid request param", err)
		return
	}
	if p.Username == "" || p.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var userP dbmodel.User
	userP.Username = p.Username
	userP.Password = p.Password

	if err := db.Create(&userP).Error; err != nil {
		println("cmn/user.go invalid create user", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": 200})
}
