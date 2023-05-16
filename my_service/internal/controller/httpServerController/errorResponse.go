package httpServerController

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
}

func NewResponseMessage(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.HTML(statusCode, "index.tmpl", map[string]string{"title": message})

}
