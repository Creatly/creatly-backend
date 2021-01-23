package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/pkg/logger"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
