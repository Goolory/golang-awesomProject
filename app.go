package main

import (
	"flag"

	"awesomeProject/config"
	"awesomeProject/dbmodel"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"awesomeProject/cmn"
	"awesomeProject/middleware"
	"awesomeProject/pro-pkg/logger"
)

func main() {

	logger.SetLevel(logger.DEBUG)
	configPath := flag.String("conf", "./config/config.json", "config file path")

	flag.Parse()

	err := config.LoadConfig(*configPath)

	if err != nil {
		logger.Fatal("config file error", err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("open db err", err)
		return
	}

	dbmodel.InitDbModel(db)

	db.LogMode(true)
	r := gin.New()

	gin.SetMode(gin.DebugMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	Api := r.Group("cmn").Use(middleware.CorsAllowHandler)
	Api.OPTIONS("/*f", middleware.CorsAllowHandler)
	Api.GET("/home", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ping"})
	})

	Api.Use(middleware.DbPrepareHandler)
	{
		Api.POST("user/register", cmn.UserRegisterHandler)

		//上传文件

		Api.POST("fileupload", cmn.FileUploadHandler)
	}

	r.NoRoute(cmn.FileServeHandler)

	r.Run()
}
