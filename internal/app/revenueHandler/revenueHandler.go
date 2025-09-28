package revenueHandler

import (
	"lab1/internal/app/revenueModel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RevenueHandler struct {
	RevenueModel *revenueModel.RevenueModel
}

func NewRevenueHandler(r *revenueModel.RevenueModel) *RevenueHandler {
	return &RevenueHandler{
		RevenueModel: r,
	}
}

func (h *RevenueHandler) GetPeriods(ctx *gin.Context) {
	var periods []revenueModel.Period
	var err error

	periodsSearchQuery := ctx.Query("periodsQuery")
	if periodsSearchQuery == "" {
		periods, err = h.RevenueModel.GetPeriods()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		periods, err = h.RevenueModel.GetPeriodsByTitle(periodsSearchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	selectedPeriodsCount, err := h.RevenueModel.GetSelectedPeriodsCount(1)
	if err != nil {
		logrus.Error(err)
		selectedPeriodsCount = 0
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"periods":              periods,
		"periodsQuery":         periodsSearchQuery,
		"selectedPeriodsCount": selectedPeriodsCount,
	})
}

func (h *RevenueHandler) GetPeriod(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	period, err := h.RevenueModel.GetPeriod(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "period.html", gin.H{
		"period": period,
	})
}

func (h *RevenueHandler) GetPeriodsApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error("error", err)
	}

	periodsApplication, err := h.RevenueModel.GetPeriodsApplication(id)
	if err != nil {
		logrus.Error("error", err)
	}

	ctx.HTML(http.StatusOK, "periodsApplication.html", gin.H{
		"periodsApplication": periodsApplication,
	})
}
