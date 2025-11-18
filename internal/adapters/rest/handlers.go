package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Util787/test-task/internal/common"
	"github.com/gin-gonic/gin"
)

type sortUsecase interface {
	SaveAndSort(ctx context.Context, num int) ([]int, error)
}

type Handler struct {
	log         *slog.Logger
	sortUsecase sortUsecase
}

type saveNumRequest struct {
	Num int `json:"num" binding:"required"`
}

func (h *Handler) saveNum(c *gin.Context) {
	log := common.LogOpAndId(c.Request.Context(), common.GetOperationName(), h.log)

	var req saveNumRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "invalid request", err)
		return
	}

	sortedArr, err := h.sortUsecase.SaveAndSort(c.Request.Context(), req.Num)
	if err != nil {
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to save number or fetch array", err)
		return
	}

	c.JSON(http.StatusOK, sortedArr)
}
