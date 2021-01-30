package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	studentCtx          = "userId"
	schoolCtx           = "school"
)

func (h *Handler) setSchoolFromRequest(c *gin.Context) {
	domainName := strings.Split(c.Request.Host, ":")[0]

	school, err := h.schoolsService.GetByDomain(c.Request.Context(), domainName)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Set(schoolCtx, school)
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

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(studentCtx, userId)
}

func getStudentId(c *gin.Context) (primitive.ObjectID, error) {
	idFromCtx, ok := c.Get(studentCtx)
	if !ok {
		return primitive.ObjectID{}, errors.New("studentCtx not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return primitive.ObjectID{}, errors.New("studentCtx is of invalid type")
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.ObjectID{}, nil
	}

	return id, nil
}
