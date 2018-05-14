package constant

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	Success      = 0
	BadParameter = 2

	UserNotFound      = 1000
	WrongPassword     = 1001
	UserAlreadyExists = 1002

	StatusInternalServerError = 500
)

var errCodeText = map[int]string{
	Success:      "Success",
	BadParameter: "Bad Parameter",

	UserNotFound:              "User Not Found",
	StatusInternalServerError: "Internal Server Error",
	WrongPassword:             "Wrong Password",
	UserAlreadyExists:         "User Already Exists",
}

func ErrCodeText(code int) string {
	return errCodeText[code]

}

func TranslaterErrCode(code int, extra ...string) string {
	var msg string
	msg = ErrCodeText(code)

	if len(extra) > 0 {
		msg = msg + ": " + strings.Join(extra, ",")
	}
	return msg

}
func ErrMsg(c *gin.Context, errCode int, extra ...string) {
	c.JSON(http.StatusOK, gin.H{"err_code": errCode, "err_msg": TranslaterErrCode(errCode, extra...)})
}
