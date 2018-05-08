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
	"awesomeProject/tool/logger"
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

	Api.Use(middleware.DbPrepareHandler)
	{
		Api.POST("/admin/login", cmn.AdminLoginHandler)
		Api.POST("/admin/register", cmn.AdminRegisterHandler)
		Api.POST("user/register", cmn.UserRegisterHandler)

		//上传文件

		Api.POST("fileupload", cmn.FileUploadHandler)
	}
	Api.Use(middleware.AdminVerifyHandler)
	{
		Api.GET("/admin/info", cmn.AdminInfoHandler)
		Api.POST("/admin/teacher/add", cmn.AdminTeacherAddHandler)
		Api.GET("/admin/teacher/list", cmn.AdminTeacherListHandler)
		Api.POST("/admin/teacher/delete", cmn.AdminTeacherDelHandler)
	}

	r.NoRoute(cmn.FileServeHandler)

	r.Run()
}
