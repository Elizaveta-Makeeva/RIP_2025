package revenueHandler

import (
	"lab2/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RevenueHandler) RegisterRevenueHandler(router *gin.Engine) {
	router.GET("/home", h.GetPeriods)
	router.GET("/period/:id", h.GetPeriod)
	router.GET("/periods-application/:id", h.GetPeriodsApplication)
	router.GET("/draft-periods-application", h.GetOrCreateDraftApplication)
	router.POST("/add-to-periods-application/:id", h.AddPeriodToApplicationForm)
	router.POST("/delete-periods-application/:id", h.DeleteApplication)
}

func (h *RevenueHandler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *RevenueHandler) errorRevenueHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *RevenueHandler) GetPeriods(ctx *gin.Context) {
	var periods []ds.Period
	var err error

	periodsSearchQuery := ctx.Query("periodsQuery")
	if periodsSearchQuery == "" {
		periods, err = h.RevenueModel.GetPeriods()
		if err != nil {
			logrus.Error(err)
			periods = []ds.Period{}
		}
	} else {
		periods, err = h.RevenueModel.GetPeriodsByTitle(periodsSearchQuery)
		if err != nil {
			logrus.Error(err)
			periods = []ds.Period{}
		}
	}

	userID := 1
	draftApp, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		logrus.Error("ошибка получения/создания черновика:", err)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"periods":              periods,
			"periodsQuery":         periodsSearchQuery,
			"selectedPeriodsCount": 0,
			"draftApplicationID":   0,
		})
		return
	}

	selectedPeriodsCount, err := h.RevenueModel.GetSelectedPeriodsCount(draftApp.ID)
	if err != nil {
		logrus.Error(err)
		selectedPeriodsCount = 0
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"periods":              periods,
		"periodsQuery":         periodsSearchQuery,
		"selectedPeriodsCount": selectedPeriodsCount,
		"draftApplicationID":   draftApp.ID,
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
		ctx.Redirect(http.StatusFound, "/home")
		return
	}

	periodsApplication, err := h.RevenueModel.GetPeriodsApplication(id)
	if err != nil {
		logrus.Error("error", err)
		ctx.Redirect(http.StatusFound, "/home")
		return
	}

	count, err := h.RevenueModel.GetSelectedPeriodsCount(id)
	if err != nil {
		logrus.Error("error", err)
		ctx.Redirect(http.StatusFound, "/home")
		return
	}

	if count == 0 {
		ctx.Redirect(http.StatusFound, "/home")
		return
	}

	ctx.HTML(http.StatusOK, "periodsApplication.html", gin.H{
		"periodsApplication": periodsApplication,
	})
}

func (h *RevenueHandler) GetOrCreateDraftApplication(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logrus.Error("error", err)
	}

	periodsApplication, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		logrus.Error("error", err)
	}

	ctx.HTML(http.StatusOK, "periodsApplication.html", gin.H{
		"periodsApplication": periodsApplication,
	})
}

func (h *RevenueHandler) AddPeriodToApplicationForm(ctx *gin.Context) {
	periodIDStr := ctx.Param("id")
	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil {
		logrus.Error("error", err)
	}

	userID := 1
	periodsApplication, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		logrus.Error("error", err)
	}

	err = h.RevenueModel.AddPeriodToApplication(periodsApplication.ID, periodID)

	redirectURL := ctx.DefaultQuery("redirect", "/home")
	ctx.Redirect(http.StatusFound, redirectURL)
}

func (h *RevenueHandler) DeleteApplication(ctx *gin.Context) {
	periodsApplicationIDStr := ctx.Param("id")
	periodsApplicationID, err := strconv.Atoi(periodsApplicationIDStr)
	if err != nil {
		logrus.Error("error", err)
	}

	err = h.RevenueModel.DeleteApplication(periodsApplicationID)
	if err != nil {
		logrus.Error("error", err)
	}

	ctx.Redirect(http.StatusFound, "/home")
}
