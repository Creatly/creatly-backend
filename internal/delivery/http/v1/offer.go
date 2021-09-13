package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
)

// @Summary Get Offer By ID
// @Tags offers
// @Description  get offer by id
// @ModuleID getOffer
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} domain.Offer
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /offers/{id} [get]
func (h *Handler) getOffer(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	offer, err := h.services.Offers.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, offer)
}
