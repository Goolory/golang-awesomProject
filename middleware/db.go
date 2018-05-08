package middleware

import (
	"net/http"

	"awesomeProject/config"
	"awesomeProject/constant"

	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func DbPrepareHandler(c *gin.Context) {
	if db == nil {
		var err error
		db, err = gorm.Open(config.GetDBName(), config.GetDBSource())
		if err != nil {
			logger.Error("Unable to connect to db (middleware db.go)", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if config.IsOrmLogEnabled() {
			db.LogMode(true)
		} else {
			db.LogMode(false)
		}
	}
	c.Set(constant.ContextDb, db)
	c.Next()

}
