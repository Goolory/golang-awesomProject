package cmn

import (
	"awesomeProject/constant"
	"awesomeProject/dbmodel"
	"awesomeProject/service"
	"awesomeProject/tool/uuid"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"awesomeProject/tool/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var idRand *rand.Rand

func init() {
	idRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}
func FileGetHandler(c *gin.Context) {
	type param struct {
		Id       uint32 `form:"id"`
		ParentId uint32 `form:"parent_id"`
		UserId   uint32 `form:"user_id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	logger.Debug(p.ParentId)
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	sql := db.Where("1=1")
	if p.Id != 0 {
		sql = sql.Where("id=?", p.Id)
	} else {
		if p.ParentId != 0 {
			sql = sql.Where("parent_id =?", p.ParentId)
		}
		if p.UserId != 0 {
			sql = sql.Where("user_id =? ", p.UserId)
		}
	}
	var files []*dbmodel.File
	if err := sql.Order("created_at desc").Find(&files).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "files": files})

}

func FileUpdateHandler(c *gin.Context) {
	type param struct {
		Id       uint32 `json:"id"`
		ParentId uint32 `json:"parent_id"`
		UserId   uint32 `json:"user_id"`
	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Info("Invalid request param", err)
		return
	}
	var file dbmodel.File
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.First(&file, p.Id).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	file.ParentId = p.ParentId
	file.UserId = p.UserId

	if err := db.Save(&file).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

func FileUploadHandler(c *gin.Context) {
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := c.Request.ParseMultipartForm(64 << 20); err != nil {
		println("maxMemory err")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("connection", "close")

	formdata := c.Request.MultipartForm

	files := formdata.File["file"]
	logger.Debug(files)

	if len(files) <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fh := files[0]
	file, err := fh.Open()

	defer file.Close()

	if err != nil {
		println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//带路径的文件名
	fileName, _ := url.QueryUnescape(fh.Filename)
	//不带路径的文件名
	fileName = filepath.Base(fileName)
	println(fileName)

	fileId, path, err := uploadFile(fileName, file, db)

	if err != nil {
		println("uploadFile err")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": 0, "data": getFullUrl(path), "id": fileId})
}

func getFullUrl(relPath string) string {
	return "http://localhost:8080" + "/res" + relPath
}

func uploadFile(filename string, f multipart.File, db *gorm.DB) (uint32, string, error) {
	tmpPath := filepath.Join("./tmp/", uuid.NewUUID().String()+filepath.Ext(filename))
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		println(err.Error())
		return 0, "", err
	}

	defer os.Remove(tmpPath)

	bytes, err := io.Copy(tmpFile, f)
	tmpFile.Close()

	if err != nil {
		println(err)
		return 0, "", err
	}

	return saveDbFile(tmpPath, filename, "", bytes, db)
}

func saveDbFile(tmpPath, filename, sourceUrl string, size int64, db *gorm.DB) (uint32, string, error) {
	var dbFile dbmodel.File

	relDir := time.Now().Format("/2006/01/02/15/0405")
	relDir = relDir + fmt.Sprintf("%03d", idRand.Intn(1000))

	if err := os.MkdirAll(filepath.Join("./file/", relDir), 0755); err != nil {
		println(err)
		return 0, "", err
	}
	dbFile.Filename = filename
	dbFile.RelPath = filepath.Join(relDir, filename)
	absPath := filepath.Join("./file/", dbFile.RelPath)
	if err := os.Rename(tmpPath, absPath); err != nil {
		println(err)
		return 0, "", err
	}
	dbFile.Filesize = size

	if err := db.Create(&dbFile).Error; err != nil {
		println(err)
		return 0, "", err
	}

	return dbFile.Id, dbFile.RelPath, nil
}

func FileServeHandler(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	relPath := strings.Split(c.Request.RequestURI, "?")[0]
	if !strings.HasPrefix(relPath, "/res") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	println(relPath)

	relPath = strings.TrimPrefix(relPath, "/res")

	relPath, _ = url.QueryUnescape(relPath)
	_, ok := service.IsThumb(relPath)

	c.Header("Connection", "close")

	absPath := filepath.Join("./file/", relPath)

	if ok {
		data, err := service.CreateThumb(absPath)
		if err != nil {
			println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, mime.TypeByExtension(filepath.Ext(filepath.Base(absPath))), data)
	} else {
		http.ServeFile(c.Writer, c.Request, absPath)
	}
}
