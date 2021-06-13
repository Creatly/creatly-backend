package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
)

// @Summary Get PromoCode By Code
// @Tags promocodes
// @Description  get promocode by code
// @ModuleID getPromo
// @Accept  json
// @Produce  json
// @Param code path string true "code"
// @Success 200 {object} domain.PromoCode
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /promocodes/{code} [get]
func (h *Handler) getPromo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "empty code param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	promocode, err := h.services.PromoCodes.GetByCode(c.Request.Context(), school.ID, code)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, promocode)
}
