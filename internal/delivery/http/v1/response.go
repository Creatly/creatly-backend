package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/pkg/logger"
)

type dataResponse struct {
	Data interface{} `json:"data"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
