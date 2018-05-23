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
	"awesomeProject/view"
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
	view.DropAllDbView(db)
	view.InitAllDbView(db)

	db.LogMode(true)
	r := gin.New()

	gin.SetMode(gin.DebugMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/static", "./file")
	r.Static("/res", "./file")

	r.StaticFile("/favicon.ico", "./image/favicon.ico")

	Api := r.Group("/cmn").Use(middleware.CorsAllowHandler)
	Api.OPTIONS("/*f", middleware.CorsAllowHandler)

	Tapi := r.Group("/cmn").Use(middleware.CorsAllowHandler)

	Sapi := r.Group("/cmn").Use(middleware.CorsAllowHandler)

	Api.Use(middleware.DbPrepareHandler)
	{
		Api.POST("/admin/login", cmn.AdminLoginHandler)
		Api.POST("/admin/register", cmn.AdminRegisterHandler)
		Api.POST("user/register", cmn.UserRegisterHandler)

		//上传文件

		Api.POST("fileupload", cmn.FileUploadHandler)
		Api.POST("fileupdate", cmn.FileUpdateHandler)
		Api.GET("getfiles", cmn.FileGetHandler)

	}
	Api.Use(middleware.AdminVerifyHandler)
	{
		Api.GET("/admin/info", cmn.AdminInfoHandler)

		Api.POST("/admin/teacher/add", cmn.AdminTeacherAddHandler)
		Api.GET("/admin/teacher/list", cmn.AdminTeacherListHandler)
		Api.POST("/admin/teacher/delete", cmn.AdminTeacherDelHandler)
		Api.GET("/admin/teacher/info", cmn.AdminTeacherInfoHandler)
		Api.POST("/admin/teacher/update", cmn.AdminTeacherUpdateHandler)
		Api.GET("/admin/teacher/all", cmn.AdminTeacherAllHandler)

		Api.POST("/admin/class/add", cmn.AdminClassAddHandler)
		Api.GET("/admin/class/list", cmn.AdminClassListHandler)
		Api.POST("/admin/class/delete", cmn.AdminClassDelHandler)
		Api.GET("/admin/class/info", cmn.AdminClassInfoHandler)
		Api.POST("/admin/class/update", cmn.AdminClassUpdateHandler)
		Api.GET("/admin/class/all", cmn.AdminClassAllHandler)

		Api.POST("/admin/student/add", cmn.AdminStudentAddHandler)
		Api.GET("/admin/student/list", cmn.AdminStudentListHandler)
		Api.POST("/admin/student/delete", cmn.AdminStudentDelHandler)
		Api.GET("/admin/student/info", cmn.AdminStudentInfoHandler)
		Api.POST("/admin/student/update", cmn.AdminStudentUpdateHandler)
		Api.GET("/admin/comment/list", cmn.AdminCommentListHandler)
		Api.POST("/admin/comment/delete", cmn.AdminCommentDelHandler)
	}

	Tapi.Use(middleware.DbPrepareHandler)
	{
		Tapi.POST("/teacher/login", cmn.TeacherLoginHandler)

	}

	Tapi.Use(middleware.TeacherVerifyHandle)
	{
		Tapi.GET("teacher/info", cmn.TeacherInfoHandler)

		Tapi.POST("/teacher/test/add", cmn.TeacherTestAddHandler)
		Tapi.GET("/teacher/test/list", cmn.TeacherTestListHandler)
		Tapi.GET("/teacher/test/info", cmn.TeacherTestInfoHandler)
		Tapi.POST("/teacher/test/delete", cmn.TeacherTestDelHandler)
		Tapi.POST("/teacher/test/update", cmn.TeacherTestUpdateHandler)
		Tapi.GET("/teacher/comment/list", cmn.TeacherCommentListHandler)
		Tapi.POST("/teacher/comment/add", cmn.TeacherCommentAddHandler)
		Tapi.GET("/teacher/homework/list", cmn.TeacherHomeworkListHandler)
		Tapi.POST("/teacher/question/add", cmn.TeacherQuestionAddHandler)
		Tapi.GET("/teacher/question/list", cmn.TeacherQuestionListHandler)
		Tapi.GET("/teacher/answer/list", cmn.TeacherAnswerListHandler)

	}

	Sapi.Use(middleware.DbPrepareHandler)
	{
		Sapi.POST("/student/login", cmn.StudentLoginHandler)

	}

	Sapi.Use(middleware.StudentVerifyHandle)
	{
		Sapi.GET("/student/info", cmn.StudentInfoHandler)
		Sapi.GET("/student/test/list", cmn.StudentTestListHandler)
		Sapi.GET("/student/test/info", cmn.StudentTestInfoHandler)
		Sapi.GET("/student/comment/list", cmn.StudentCommentListHandler)
		Sapi.POST("/student/comment/add", cmn.StudentCommentAddHandler)
		Sapi.POST("/student/homework/add", cmn.StudentHomeworkAddHandler)
		Sapi.POST("/student/answer/add", cmn.StudentAnswerAddHandler)
		Sapi.GET("/student/question/list", cmn.StudentQuestionListHandler)

	}

	r.NoRoute(cmn.FileServeHandler)

	r.Run()
}
