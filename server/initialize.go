package server

import "github.com/gin-gonic/gin"

func Initialize(rGroup *gin.RouterGroup) {
	Mount(rGroup)
}
