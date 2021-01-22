package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"net/http"
	"strings"
)

const (
	schoolCtx = "school"
)

// TODO extract domain from host logic
func (h *Handler) setSchoolFromRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		domainName := strings.Split(c.Request.Host, ":")[0]

		school, err := h.schoolsService.GetByDomain(c.Request.Context(), domainName)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set(schoolCtx, school)
	}
}

func getSchoolFromContext(c *gin.Context) (domain.School, error) {
	value, ex := c.Get(schoolCtx)
	if !ex {
		return domain.School{}, errors.New("school is missing from ctx")
	}

	school, ok := value.(domain.School)
	if !ok {
		return domain.School{}, errors.New("failed to convert value from ctx to domain.School")
	}

	return school, nil
}
