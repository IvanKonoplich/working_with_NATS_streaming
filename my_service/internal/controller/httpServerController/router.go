package httpServerController

import "github.com/gin-gonic/gin"

func (c *HttpServerController) InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*")
	router.GET("/:id", c.HandleGetById)
	return router
}
